package loadbalance

import (
	"context"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"net/http"
	"sync"
	"time"
)

// check interval
const checkInterval = 5

type LoadBalance interface {
	Close()
	SetMode(mode int)
	GetChainId() uint64
	NextClient() *rpcclient.EvmClient
}

type loadBalance struct {
	chainId       uint64
	urls          []string
	healthyCliMap map[int]*rpcclient.EvmClient
	currentIndex  int
	lock          *sync.RWMutex
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
		chainId:       cli.GetChainId(),
		urls:          urls,
		healthyCliMap: make(map[int]*rpcclient.EvmClient, len(urls)),
		currentIndex:  0,
		lock:          &sync.RWMutex{},
		ctx:           ctx,
		cancel:        cancel,
	}
	l.healthyCliMap[0] = cli
	go l.checkHealth()
	return l
}

func connClient(url string) *rpcclient.EvmClient {
	client, _ := rpcclient.NewEvmRpcClient(url)
	return client
}

// check health
func (l *loadBalance) checkHealth() {
	ticker := time.NewTicker(checkInterval * time.Second)

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			for i, url := range l.urls {
				if !isHealthy(url) {
					l.lock.Lock()
					if _, ok := l.healthyCliMap[i]; ok {
						l.healthyCliMap[i].Close()
						delete(l.healthyCliMap, i)
					}
					l.lock.Unlock()
					continue
				}

				l.lock.Lock()
				if _, ok := l.healthyCliMap[i]; !ok {
					cli := connClient(url)
					if cli == nil || l.chainId != cli.GetChainId() { // check chain id
						l.lock.Unlock()
						continue
					}
					l.healthyCliMap[i] = cli
				}
				l.lock.Unlock()
			}
		}

	}
}

// check url is healthy
func isHealthy(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	_, err = http.DefaultClient.Do(req)
	//if err == nil {
	//	fmt.Printf("healthy: %v\n", url)
	//}
	return err == nil
}

// set mode
func (l *loadBalance) SetMode(mode int) {
	l.lock.Lock()
	defer l.lock.Unlock()
}

// get client
func (l *loadBalance) NextClient() *rpcclient.EvmClient {
	l.lock.RLock()
	if len(l.healthyCliMap) == 0 {
		// All clients are down
		l.lock.RUnlock()
		return nil
	}
	l.lock.RUnlock()

	// get next client
	l.lock.Lock()
	defer l.lock.Unlock()
	l.currentIndex = (l.currentIndex + 1) % len(l.healthyCliMap)
	return l.healthyCliMap[l.currentIndex]
}

func (l *loadBalance) GetChainId() uint64 {
	return l.chainId
}

func (l *loadBalance) Close() {
	l.cancel()
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, cli := range l.healthyCliMap {
		cli.Close()
	}
}
