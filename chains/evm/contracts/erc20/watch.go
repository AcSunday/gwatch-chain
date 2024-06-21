package erc20

import (
	"github.com/ethereum/go-ethereum/common"
	"gwatch_chain/chains/evm/contracts/abs"
)

type ERC20 struct {
	abs.Contract
}

func New(addr common.Address, attrs *abs.Attrs) *ERC20 {
	e := &ERC20{
		Contract: abs.Contract{
			Addr: addr,
		},
	}
	e.Init(*attrs)
	e.IsRunning.Store(true)
	return e
}
