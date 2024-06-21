package gwatch_chain

import (
	"context"
	"errors"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// quick start

type IWatch interface {
	Watch() error

	Close() error
	RegisterWatchEvent(topics ...abs.Event) error
	RegisterEventHook(event abs.Event, f func(log types.Log) error) error
	UpdateProcessedBlockNumber(num uint64) error
	GetProcessedBlockNumber() uint64
}

type Options struct {
	abs.Attrs
}

type watch struct {
	client *rpcclient.EvmClient
	abs.IContract
}

func (w *watch) Watch() error {
	return w.IContract.Watch(w.client)
}

func (w *watch) Close() error {
	w.IContract.Close()
	w.client.Close()
	return nil
}

func NewERC721Watch(rawurl string, addr common.Address, ops *Options) (IWatch, error) {
	client := rpcclient.MustNewEvmRpcClient(rawurl)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, errors.New("get chain id err:" + err.Error())
	}

	e := erc721.New(addr, &ops.Attrs)
	e.ChainId = chainID.Uint64()

	return &watch{client: client, IContract: e}, nil
}
