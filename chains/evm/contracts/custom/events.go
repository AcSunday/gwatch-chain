package custom

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/ethereum/go-ethereum/crypto"
)

// CustomEvent
//
//	Transfer string case: Transfer(address,address,uint256)
func CustomEvent(fStr string) abs.Event {
	return abs.Event(crypto.Keccak256Hash([]byte(fStr)).Hex())
}
