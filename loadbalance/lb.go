package loadbalance

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// check interval
	checkInterval          = 5
	delayedClosingInterval = 20

	unhealthyTolerateVal = 3
	maxIdleConns         = 50
	idleConnTimeout      = 90 * time.Second
	clientTimeout        = 1 * time.Second
)

// RPCClient 定义了RPC客户端需要实现的通用接口
type RPCClient interface {
	comparable
	// Close 关闭客户端连接
	Close()
	// GetRawUrl 获取原始URL
	GetRawUrl() string
	// GetChainId 获取链ID
	GetChainId() uint64
}

// ClientFactory 定义了创建客户端的工厂函数类型
type ClientFactory[T RPCClient] func(url string) (T, error)

// 客户端连接池
var clientPool = sync.Pool{
	New: func() interface{} {
		return &http.Client{
			Timeout: clientTimeout,
			Transport: &http.Transport{
				MaxIdleConns:       maxIdleConns,
				IdleConnTimeout:    idleConnTimeout,
				DisableCompression: true,
			},
		}
	},
}

// LoadBalance 负载均衡器接口
type LoadBalance[T RPCClient] interface {
	Close()
	SetMode(mode int)
	GetChainId() uint64
	NextClient() T
	ReleaseClient(T)
}

// nodeInfo 节点信息
type nodeInfo[T RPCClient] struct {
	client       T
	unhealthyCnt atomic.Int32
	lastCheck    atomic.Int64 // Unix timestamp
	refCount     atomic.Int32
}

// loadBalance 负载均衡器实现
type loadBalance[T RPCClient] struct {
	chainId      uint64
	urls         []string
	nodes        []*nodeInfo[T]
	nodeSnapshot atomic.Value  // []*nodeInfo[T]
	currentIndex atomic.Uint32 // 使用 Uint32 避免负数
	ctx          context.Context
	cancel       context.CancelFunc
	factory      ClientFactory[T]
}

// New 创建新的负载均衡器
func New[T RPCClient](urls []string, factory ClientFactory[T]) LoadBalance[T] {
	if len(urls) == 0 {
		return nil
	}

	// 创建第一个客户端
	cli, err := factory(urls[0])
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	l := &loadBalance[T]{
		chainId: cli.GetChainId(),
		urls:    urls,
		nodes:   make([]*nodeInfo[T], len(urls)),
		ctx:     ctx,
		cancel:  cancel,
		factory: factory,
	}

	// 初始化节点信息
	l.nodes[0] = &nodeInfo[T]{client: cli}

	for i := 1; i < len(urls); i++ {
		if client, err := factory(urls[i]); err == nil {
			l.nodes[i] = &nodeInfo[T]{client: client}
		}
	}

	// 初始化快照
	l.updateNodesSnapshot()

	// 启动健康检查
	go l.healthCheckLoop()

	return l
}

func (l *loadBalance[T]) updateNodesSnapshot() {
	healthyNodes := make([]*nodeInfo[T], 0, len(l.nodes))
	for _, node := range l.nodes {
		if node != nil && !l.isZero(node.client) && node.unhealthyCnt.Load() < unhealthyTolerateVal {
			healthyNodes = append(healthyNodes, node)
		}
	}
	if len(healthyNodes) > 0 {
		l.nodeSnapshot.Store(healthyNodes)
	}
}

// isZero 检查客户端是否为零值
func (l *loadBalance[T]) isZero(client T) bool {
	var zero T
	return client == zero
}

func (l *loadBalance[T]) healthCheckLoop() {
	ticker := time.NewTicker(checkInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			l.parallelHealthCheck()
		}
	}
}

func (l *loadBalance[T]) parallelHealthCheck() {
	var wg sync.WaitGroup
	wg.Add(len(l.urls))

	for i := range l.urls {
		go func(idx int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("health check panic recovered: %v", r)
				}
			}()

			node := l.nodes[idx]
			if node == nil || l.isZero(node.client) {
				return
			}

			now := time.Now().Unix()
			// 添加健康检查间隔
			if now-node.lastCheck.Load() < checkInterval {
				return
			}

			if !isHealthy(l.urls[idx]) {
				unhealthyCnt := node.unhealthyCnt.Add(1)
				if unhealthyCnt >= unhealthyTolerateVal {
					oldClient := node.client
					var zero T
					node.client = zero
					if node.refCount.Load() == 0 {
						oldClient.Close()
					} else {
						go l.delayedClosing(oldClient)
					}
				}
			} else {
				// 健康时重置不健康计数
				node.unhealthyCnt.Store(0)
				node.lastCheck.Store(now)

				if l.isZero(node.client) {
					if newClient, err := l.factory(l.urls[idx]); err == nil && l.chainId == newClient.GetChainId() {
						node.client = newClient
					}
				}
			}
		}(i)
	}

	wg.Wait()
	l.updateNodesSnapshot()
}

func (l *loadBalance[T]) delayedClosing(cli T) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("delayed closing panic recovered: %v", r)
		}
	}()

	ctx, cancel := context.WithTimeout(l.ctx, 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			cli.Close()
			return
		case <-ticker.C:
			found := false
			for _, node := range l.nodes {
				if node != nil && !l.isZero(node.client) && node.client == cli {
					if node.refCount.Load() == 0 {
						cli.Close()
						return
					}
					found = true
					break
				}
			}
			if !found {
				return
			}
		}
	}
}

func isHealthy(url string) bool {
	client := clientPool.Get().(*http.Client)
	defer clientPool.Put(client)

	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (l *loadBalance[T]) SetMode(mode int) {
	// 预留接口，支持不同的负载均衡模式
}

func (l *loadBalance[T]) NextClient() T {
	nodesI := l.nodeSnapshot.Load()
	if nodesI == nil {
		var zero T
		return zero
	}

	nodes := nodesI.([]*nodeInfo[T])
	if len(nodes) == 0 {
		var zero T
		return zero
	}

	// 使用 Uint32，永远为正数
	idx := l.currentIndex.Add(1)
	node := nodes[idx%uint32(len(nodes))]
	node.refCount.Add(1)
	return node.client
}

func (l *loadBalance[T]) ReleaseClient(cli T) {
	for _, node := range l.nodes {
		if node != nil && !l.isZero(node.client) && node.client == cli {
			node.refCount.Add(-1)
			return
		}
	}
}

func (l *loadBalance[T]) GetChainId() uint64 {
	return l.chainId
}

func (l *loadBalance[T]) Close() {
	l.cancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for _, node := range l.nodes {
		if node != nil && !l.isZero(node.client) {
			wg.Add(1)
			go func(n *nodeInfo[T]) {
				defer wg.Done()

				ticker := time.NewTicker(100 * time.Millisecond)
				defer ticker.Stop()

				cli := n.client
				for {
					select {
					case <-ctx.Done():
						cli.Close()
						return
					case <-ticker.C:
						if n.refCount.Load() == 0 {
							cli.Close()
							return
						}
					}
				}
			}(node)
		}
	}

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-ctx.Done():
		// log.Printf("load balancer close timeout after 10s")
	case <-waitCh:
		// log.Printf("load balancer closed gracefully")
	}
}
