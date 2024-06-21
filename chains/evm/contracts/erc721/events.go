package erc721

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/ethereum/go-ethereum/crypto"
)

// TransferEvent
//
//	Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
func TransferEvent() abs.Event {
	return abs.Event(crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex())
}

// ApprovalEvent
//
// Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
func ApprovalEvent() abs.Event {
	return abs.Event(crypto.Keccak256Hash([]byte("Approval(address,address,uint256)")).Hex())
}

// ApprovalForAllEvent
//
// ApprovalForAll(address indexed owner, address indexed operator, bool approved);
func ApprovalForAllEvent() abs.Event {
	return abs.Event(crypto.Keccak256Hash([]byte("ApprovalForAll(address,address,bool)")).Hex())
}
