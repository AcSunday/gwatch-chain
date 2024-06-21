package contracts

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc20"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"testing"
)

func TestGetEvents(t *testing.T) {
	t.Logf("--- ERC20 TransferEvent: %s", erc20.TransferEvent())
	t.Logf("--- ERC20 ApprovalEvent: %s", erc20.ApprovalEvent())

	t.Logf("--- ERC721 TransferEvent: %s", erc721.TransferEvent())
	t.Logf("--- ERC721 ApprovalEvent: %s", erc721.ApprovalEvent())
	t.Logf("--- ERC721 ApprovalForAllEvent: %s", erc721.ApprovalForAllEvent())
}
