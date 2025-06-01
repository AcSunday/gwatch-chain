package loadbalance

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AcSunday/gwatch-chain/rpcclient"
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

type LoadBalance interface {
	Close()
	SetMode(mode int)
	GetChainId() uint64
	NextClient() *rpcclient.EvmClient
}

type nodeInfo struct {
	client       *rpcclient.EvmClient
	unhealthyCnt int32
	lastCheck    int64 // Unix timestamp
}

type loadBalance struct {
	chainId       uint64
	urls          []string
	nodes         []nodeInfo
	nodesSnapshot atomic.Value // []*nodeInfo
	currentIndex  atomic.Int32
	ctx           context.Context
	cancel        context.CancelFunc
}

func New(urls []string) LoadBalance {
	if len(urls) == 0 {
		return nil
	}
	cli := rpcclient.MustNewEvmRpcClient(urls[0])

	ctx, cancel := context.WithCancel(context.Background())

	l := &loadBalance{
		chainId: cli.GetChainId(),
		urls:    urls,
		nodes:   make([]nodeInfo, len(urls)),
		ctx:     ctx,
		cancel:  cancel,
	}

	// 初始化节点信息
	l.nodes[0] = nodeInfo{client: cli}
	for i := 1; i < len(urls); i++ {
		if client := connClient(urls[i]); client != nil {
			l.nodes[i] = nodeInfo{client: client}
		}
	}

	// 初始化快照
	l.updateNodesSnapshot()

	// 启动健康检查
	go l.healthCheckLoop()

	return l
}

func (l *loadBalance) updateNodesSnapshot() {
	healthyNodes := make([]*nodeInfo, 0, len(l.nodes))
	for i := range l.nodes {
		if l.nodes[i].client != nil && atomic.LoadInt32(&l.nodes[i].unhealthyCnt) < unhealthyTolerateVal {
			healthyNodes = append(healthyNodes, &l.nodes[i])
		}
	}
	l.nodesSnapshot.Store(healthyNodes)
}

func connClient(url string) *rpcclient.EvmClient {
	client, _ := rpcclient.NewEvmRpcClient(url)
	return client
}

func (l *loadBalance) healthCheckLoop() {
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

func (l *loadBalance) parallelHealthCheck() {
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

			node := &l.nodes[idx]
			if node.client == nil {
				return
			}

			if !isHealthy(l.urls[idx]) {
				unhealthyCnt := atomic.AddInt32(&node.unhealthyCnt, 1)
				if unhealthyCnt >= unhealthyTolerateVal {
					oldClient := node.client
					node.client = nil
					go l.delayedClosing(oldClient)
				}
			} else {
				atomic.StoreInt32(&node.unhealthyCnt, 0)
				atomic.StoreInt64(&node.lastCheck, time.Now().Unix())

				if node.client == nil {
					if newClient := connClient(l.urls[idx]); newClient != nil && l.chainId == newClient.GetChainId() {
						node.client = newClient
					}
				}
			}
		}(i)
	}

	wg.Wait()
	l.updateNodesSnapshot()
}

func (l *loadBalance) delayedClosing(cli *rpcclient.EvmClient) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("delayed closing panic recovered: %v", r)
		}
	}()
	time.Sleep(delayedClosingInterval * time.Second)
	cli.Close()
}

// check url is healthy
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

// set mode
func (l *loadBalance) SetMode(mode int) {
	// 预留接口，支持不同的负载均衡模式
}

// get client
func (l *loadBalance) NextClient() *rpcclient.EvmClient {
	nodes := l.nodesSnapshot.Load().([]*nodeInfo)
	if len(nodes) == 0 {
		return nil
	}

	// 使用原子操作更新索引
	idx := l.currentIndex.Add(1)
	node := nodes[idx%int32(len(nodes))]

	return node.client
}

func (l *loadBalance) GetChainId() uint64 {
	return l.chainId
}

func (l *loadBalance) Close() {
	l.cancel()
	for i := range l.nodes {
		if l.nodes[i].client != nil {
			l.nodes[i].client.Close()
		}
	}
}
