package erc721

import (
	"github.com/ethereum/go-ethereum/common"
	"gwatch_chain/chains/evm/contracts/abs"
)

type ERC721 struct {
	abs.Contract
}

func New(addr common.Address, attrs *abs.Attrs) *ERC721 {
	e := &ERC721{
		Contract: abs.Contract{
			Addr: addr,
		},
	}
	e.Init(*attrs)
	e.IsRunning.Store(true)
	e.Decimals = 0
	return e
}
