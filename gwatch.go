package gwatch

import (
	"errors"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc20"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"github.com/AcSunday/gwatch-chain/loadbalance"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// quick start

type IWatch interface {
	Watch() error

	Close() error
	DoneSignal() <-chan struct{}
	RegisterWatchEvent(events ...abs.Event) error
	RegisterWatchTopics(topicsIndex int, topics ...common.Hash) error
	RegisterEventHook(event abs.Event, f func(client *rpcclient.EvmClient, log types.Log) error) error
	UpdateProcessedBlockNumber(num uint64) error
	GetProcessedBlockNumber() uint64
}

type Options struct {
	abs.Attrs
}

type watch struct {
	lb loadbalance.LoadBalance
	abs.IContract
}

func (w *watch) Watch() error {
	cli := w.lb.NextClient()
	for i := 0; i < 3; i++ {
		if cli != nil {
			break
		}
		cli = w.lb.NextClient()
	}
	if cli == nil {
		return errors.New("no clients available, failed to connect to blockchain")
	}
	return w.IContract.Watch(cli)
}

func (w *watch) Close() error {
	w.IContract.Close()
	w.lb.Close()
	return nil
}

func NewERC721Watch(rawurls []string, addrs []common.Address, ops *Options) (IWatch, error) {
	l := loadbalance.New(rawurls)

	e := erc721.New(addrs, &ops.Attrs)
	e.ChainId = l.GetChainId()

	return &watch{lb: l, IContract: e}, nil
}

func NewGeneralWatch(rawurls []string, addrs []common.Address, ops *Options) (IWatch, error) {
	l := loadbalance.New(rawurls)

	e := erc20.New(addrs, &ops.Attrs)
	e.ChainId = l.GetChainId()

	return &watch{lb: l, IContract: e}, nil
}

func NewLoadBalanceGeneralWatch(lb loadbalance.LoadBalance, addrs []common.Address, ops *Options) (IWatch, error) {
	e := erc20.New(addrs, &ops.Attrs)
	e.ChainId = lb.GetChainId()

	return &watch{lb: lb, IContract: e}, nil
}
