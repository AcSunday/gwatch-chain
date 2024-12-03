package erc20

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/ethereum/go-ethereum/common"
)

type ERC20 struct {
	abs.Contract
}

func New(addrs []common.Address, attrs *abs.Attrs) *ERC20 {
	e := &ERC20{
		Contract: abs.Contract{
			Addrs: addrs,
		},
	}
	e.Init(*attrs)
	e.IsRunning.Store(true)
	return e
}
