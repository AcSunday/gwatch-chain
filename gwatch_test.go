package gwatch

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"github.com/AcSunday/gwatch-chain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
)

func TestQuickStart(t *testing.T) {
	e, err := NewERC721Watch(
		[]string{"https://sepolia.infura.io/v3/f6ef0da20fa14730ae77a316d88c0516"},
		common.HexToAddress("0x9643E463b77a6c562eb6d459980622fbB8a91e1D"),
		&Options{
			Attrs: abs.Attrs{
				Chain:                "sepolia",
				Name:                 "POPTAG",
				Symbol:               "POP",
				DeployedBlockNumber:  6060135,
				ProcessedBlockNumber: 6152533,
				WatchBlockLimit:      50,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	e.RegisterWatchEvent(erc721.ApprovalEvent(), erc721.TransferEvent(), erc721.ApprovalForAllEvent())
	e.RegisterEventHook(erc721.TransferEvent(), func(client *ethclient.Client, log types.Log) error {
		t.Logf("------ transfer log txhash: %s ------", log.TxHash)
		from := utils.EventAddressHashFormat(log.Topics[1])
		to := utils.EventAddressHashFormat(log.Topics[2])
		t.Logf("From: %s to: %s tokenID: %s", from, to, log.Topics[3].Big().String())
		return nil
	})

	for i := 0; i < 10; i++ {
		err = e.Watch()
		if err != nil {
			t.Fatal(err)
		}
	}
}
