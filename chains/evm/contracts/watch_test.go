package contracts

import (
	"context"
	"fmt"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc20"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/AcSunday/gwatch-chain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"testing"
)

const (
	rawurl = "https://sepolia.infura.io/v3/f6ef0da20fa14730ae77a316d88c0516"

	ERC20ContractAddr = "0x1c7d4b196cb0c7b01d743fbc6116a902379c7238"
	NFTContractAddr   = "0x9643E463b77a6c562eb6d459980622fbB8a91e1D"
)

func TestWatchERC20(t *testing.T) {
	client := rpcclient.MustNewEvmRpcClient(rawurl)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	e := erc20.New([]common.Address{common.HexToAddress(ERC20ContractAddr)}, &abs.Attrs{
		ChainId:              chainID.Uint64(),
		Chain:                "sepolia",
		Name:                 "USDC",
		Symbol:               "USDC",
		Decimals:             6,
		DeployedBlockNumber:  4848135,
		ProcessedBlockNumber: 6136772,
		WatchBlockLimit:      2,
	})

	e.RegisterWatchEvent(erc20.ApprovalEvent(), erc20.TransferEvent())
	e.RegisterWatchTopics(2, common.HexToHash("0x47ccd6b8e3e0e1b84ad818842fd68b209a6a9cd7"))
	e.RegisterEventHook(erc20.TransferEvent(), func(client *rpcclient.EvmClient, log types.Log) error {
		t.Logf("------ transfer log txhash: %s ------", log.TxHash)
		from := fmt.Sprintf("0x%s", log.Topics[1].String()[26:])
		to := fmt.Sprintf("0x%s", log.Topics[2].String()[26:])
		amount := common.BytesToHash(log.Data)
		//t.Logf("From: %s to: %s amount: %s", from, to, amount.Big().String())
		t.Logf("From: %s to: %s amount: %f", from, to, utils.AmountToStr(int64(e.Decimals), amount.Big()))
		return nil
	})

	for i := 0; i < 1; i++ {
		err = e.Watch(client)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestWatchERC721(t *testing.T) {
	client := rpcclient.MustNewEvmRpcClient(rawurl)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	e := erc721.New([]common.Address{common.HexToAddress(NFTContractAddr)}, &abs.Attrs{
		ChainId:              chainID.Uint64(),
		Chain:                "sepolia",
		Name:                 "POPTAG",
		Symbol:               "POP",
		DeployedBlockNumber:  6060135,
		ProcessedBlockNumber: 6152533,
		WatchBlockLimit:      50,
	})

	e.RegisterWatchEvent(erc721.ApprovalEvent(), erc721.TransferEvent(), erc721.ApprovalForAllEvent())
	e.RegisterEventHook(erc721.TransferEvent(), func(client *rpcclient.EvmClient, log types.Log) error {
		t.Logf("------ transfer log txhash: %s ------", log.TxHash)
		from := fmt.Sprintf("0x%s", log.Topics[1].String()[26:])
		to := fmt.Sprintf("0x%s", log.Topics[2].String()[26:])
		t.Logf("From: %s to: %s tokenID: %s", from, to, log.Topics[3].Big().String())
		return nil
	})

	for i := 0; i < 10; i++ {
		err = e.Watch(client)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestWatchERC1155(t *testing.T) {
	// TODO
}
