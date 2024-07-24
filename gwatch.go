package gwatch

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc20"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"github.com/AcSunday/gwatch-chain/loadbalance"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// quick start

type IWatch interface {
	Watch() error

	Close() error
	DoneSignal() <-chan struct{}
	RegisterWatchEvent(topics ...abs.Event) error
	RegisterEventHook(event abs.Event, f func(client *ethclient.Client, log types.Log) error) error
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
	//for {
	//	if cli != nil {
	//		break
	//	}
	//	cli = w.lb.NextClient()
	//}
	return w.IContract.Watch(cli)
}

func (w *watch) Close() error {
	w.IContract.Close()
	w.lb.Close()
	return nil
}

func NewERC721Watch(rawurls []string, addr common.Address, ops *Options) (IWatch, error) {
	l := loadbalance.New(rawurls)

	e := erc721.New(addr, &ops.Attrs)
	e.ChainId = l.GetChainId()

	return &watch{lb: l, IContract: e}, nil
}

func NewGeneralWatch(rawurls []string, addr common.Address, ops *Options) (IWatch, error) {
	l := loadbalance.New(rawurls)

	e := erc20.New(addr, &ops.Attrs)
	e.ChainId = l.GetChainId()

	return &watch{lb: l, IContract: e}, nil
}
