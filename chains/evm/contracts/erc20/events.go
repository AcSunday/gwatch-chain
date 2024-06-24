package erc20

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/ethereum/go-ethereum/crypto"
)

// string event
const (
	transferEvent = "Transfer"
	approvalEvent = "Approval"
)

// TransferEvent
//
// Transfer(address indexed from, address indexed to, uint256 value);
func TransferEvent() abs.Event {
	return abs.Event(crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex())
}

// ApprovalEvent
//
// Approval(address indexed owner, address indexed spender, uint256 value);
func ApprovalEvent() abs.Event {
	return abs.Event(crypto.Keccak256Hash([]byte("Approval(address,address,uint256)")).Hex())
}

// EventToName ...
func EventToName(event abs.Event) string {
	switch event {
	case TransferEvent():
		return transferEvent
	case ApprovalEvent():
		return approvalEvent
	}

	return ""
}
