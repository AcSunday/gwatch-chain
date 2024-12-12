package abs

import (
	"context"
	"errors"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"sync"
	"sync/atomic"
)

type Event string

const DefaultWatchLimit = 20

type Attrs struct {
	ChainId              uint64
	Chain                string
	Name                 string
	Symbol               string
	Decimals             uint8
	DeployedBlockNumber  uint64 // contract deployment height
	ProcessedBlockNumber uint64 // has been processed on block number, default is DeployedBlockNumber
	WatchBlockLimit      int64  // Limit the number of blocks scanned each time, default is 20
}

type Contract struct {
	Attrs

	Addrs  []common.Address
	Topics [][]common.Hash

	IsRunning atomic.Bool
	IsClose   atomic.Bool

	handleFunc map[Event]func(client *rpcclient.EvmClient, log types.Log) error // key is event
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func (c *Contract) Init(attrs Attrs) {
	c.Topics = make([][]common.Hash, 1)
	c.handleFunc = make(map[Event]func(client *rpcclient.EvmClient, log types.Log) error, 4)

	c.Attrs = attrs
	if c.WatchBlockLimit <= 0 {
		c.WatchBlockLimit = DefaultWatchLimit
	}
	if c.DeployedBlockNumber-1 > 0 && c.ProcessedBlockNumber < c.DeployedBlockNumber {
		c.ProcessedBlockNumber = c.DeployedBlockNumber - 1
	}

	//c.IsRunning.Store(false)
	c.IsClose.Store(false)
	c.mu = sync.RWMutex{}
	c.ctx, c.cancel = context.WithCancel(context.Background())
}

func (c *Contract) Close() error {
	if !c.IsClose.Load() {
		c.cancel()
		c.IsRunning.Store(false)
		c.IsClose.Store(true)
	}
	return nil
}

func (c *Contract) DoneSignal() <-chan struct{} {
	if c.IsRunning.Load() && !c.IsClose.Load() {
		return c.ctx.Done()
	}
	return nil
}

// RegisterWatchEvent topics[0] is smart contract event
func (c *Contract) RegisterWatchEvent(events ...Event) error {
	if c.IsRunning.Load() && len(c.Topics[0]) > 0 {
		return errors.New("already running, Registration of events is prohibited")
	}

	sli := make([]common.Hash, 0, len(events))
	for _, event := range events {
		sli = append(sli, common.HexToHash(event.String()))
	}
	c.mu.Lock()
	c.Topics[0] = append(c.Topics[0], sli...)
	c.mu.Unlock()
	return nil
}

// RegisterWatchTopics topic is smart contract event parameter
//
//	topicsIndex: [0-3]
//	event topics
func (c *Contract) RegisterWatchTopics(topicsIndex int, topics ...common.Hash) error {
	if c.IsClose.Load() {
		return errors.New("already closed, Registration of topics is prohibited")
	}
	if topicsIndex > 3 {
		return errors.New("invalid topic index")
	}

	c.mu.Lock()
	for i := len(c.Topics); i <= topicsIndex; i++ {
		c.Topics = append(c.Topics, []common.Hash{})
	}
	if topicsIndex < len(c.Topics) {
		c.Topics[topicsIndex] = append(c.Topics[topicsIndex], topics...)
	}
	c.mu.Unlock()
	return nil
}

// RegisterEventHook Hook is a function that handles event,
// HandleEvent method call this Hook
func (c *Contract) RegisterEventHook(event Event, f func(client *rpcclient.EvmClient, log types.Log) error) error {
	if c.IsClose.Load() {
		return errors.New("already closed, Registration of event hook is prohibited")
	}
	c.mu.Lock()
	c.handleFunc[event] = f
	c.mu.Unlock()
	return nil
}

// HandleEvent method call Hook
func (c *Contract) HandleEvent(client *rpcclient.EvmClient, event Event, log types.Log) error {
	if !c.IsRunning.Load() {
		return errors.New("not running, handle event is prohibited")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if f, ok := c.handleFunc[event]; ok {
		return f(client, log)
	}
	return nil
}

// UpdateProcessedBlockNumber ...
func (c *Contract) UpdateProcessedBlockNumber(num uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ProcessedBlockNumber = num
	return nil
}

// UpdateProcessedBlockNumber ...
func (c *Contract) GetProcessedBlockNumber() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ProcessedBlockNumber
}

func (e Event) String() string {
	return string(e)
}
