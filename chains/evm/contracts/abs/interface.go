package abs

import (
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type IContract interface {
	Init(attrs Attrs)
	Close() error
	DoneSignal() <-chan struct{}
	RegisterWatchEvent(events ...Event) error
	RegisterWatchTopics(topicsIndex int, topics ...common.Hash) error
	RegisterEventHook(event Event, f func(client *rpcclient.EvmClient, log types.Log) error) error
	HandleEvent(client *rpcclient.EvmClient, event Event, log types.Log) error
	UpdateProcessedBlockNumber(num uint64) error
	GetProcessedBlockNumber() uint64
	Scan(client *rpcclient.EvmClient) error
	GetContractDesc(addr string) (ContractDesc, error)
}
