package abs

import (
	"github.com/ethereum/go-ethereum/core/types"
	"gwatch_chain/rpcclient"
)

type IContract interface {
	Init(attrs Attrs)
	Close() error
	DoneSignal() <-chan struct{}
	RegisterWatchEvent(topics ...Event) error
	RegisterEventHook(event Event, f func(log types.Log) error) error
	HandleEvent(event Event, log types.Log) error
	UpdateProcessedBlockNumber(num uint64) error
	GetProcessedBlockNumber() uint64
	Watch(client *rpcclient.EvmClient) error
}
