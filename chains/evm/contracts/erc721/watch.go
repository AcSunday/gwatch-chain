package erc721

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/ethereum/go-ethereum/common"
)

type ERC721 struct {
	abs.Contract
}

func New(addrs []common.Address, attrs *abs.Attrs) *ERC721 {
	e := &ERC721{
		Contract: abs.Contract{
			Addrs: addrs,
		},
	}
	e.Init(*attrs)
	return e
}
