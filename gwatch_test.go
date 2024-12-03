package gwatch

import (
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc20"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"github.com/AcSunday/gwatch-chain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
)

func TestQuickStartERC20(t *testing.T) {
	e, err := NewGeneralWatch(
		[]string{"https://1rpc.io/bnb"},
		[]common.Address{common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")}, // contract address
		&Options{
			Attrs: abs.Attrs{
				Chain:                "BNB",
				Name:                 "USDT",
				Symbol:               "USDT",
				Decimals:             18,
				DeployedBlockNumber:  43512850,
				ProcessedBlockNumber: 43512850,
				WatchBlockLimit:      2,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	e.RegisterWatchEvent(erc20.ApprovalEvent(), erc20.TransferEvent())
	//e.RegisterWatchTopics(1, common.HexToHash("0xD7D761ce2e145FF4b72321F2075679ea42255286")) // filter from addr
	e.RegisterWatchTopics(2, common.HexToHash("0x9560d82d93a6a6d204df56f101964b26ce61e999")) // filter to addr
	e.RegisterEventHook(erc20.TransferEvent(), func(client *ethclient.Client, log types.Log) error {
		t.Logf("------ %d transfer log txhash: %s ------", log.BlockNumber, log.TxHash)
		from := utils.EventAddressHashFormat(log.Topics[1])
		to := utils.EventAddressHashFormat(log.Topics[2])
		amount := common.BytesToHash(log.Data)
		//t.Logf("From: %s to: %s amount: %s", from, to, amount.Big().String())
		t.Logf("From: %s to: %s amount: %f", from, to, utils.AmountToStr(18, amount.Big()))
		return nil
	})

	for i := 0; i < 10; i++ {
		err = e.Watch()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestQuickStartERC721(t *testing.T) {
	e, err := NewERC721Watch(
		[]string{"https://sepolia.infura.io/v3/f6ef0da20fa14730ae77a316d88c0516"},
		[]common.Address{common.HexToAddress("0x9643E463b77a6c562eb6d459980622fbB8a91e1D")},
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
		t.Logf("------ %d transfer log txhash: %s ------", log.BlockNumber, log.TxHash)
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
