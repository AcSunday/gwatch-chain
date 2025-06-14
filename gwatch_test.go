package gwatch

import (
	"testing"

	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/abs"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc20"
	"github.com/AcSunday/gwatch-chain/chains/evm/contracts/erc721"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/AcSunday/gwatch-chain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	ERC20  = "0x55d398326f99059ff775485246999027b3197955"
	ERC721 = "0x9643E463b77a6c562eb6d459980622fbB8a91e1D"
)

func TestQuickStartERC20(t *testing.T) {
	e, err := NewGeneralWatch(
		[]string{"https://1rpc.io/bnb"},
		[]common.Address{common.HexToAddress(ERC20)}, // contract address
		&Options{
			Attrs: abs.Attrs{
				Chain:                "BNB",
				DeployedBlockNumber:  43512850,
				ProcessedBlockNumber: 43512850,
				WatchBlockLimit:      2,
				ContractToDesc: map[string]abs.ContractDesc{
					ERC20: {
						Name:     "USDT",
						Symbol:   "USDT",
						Decimals: 18,
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	e.RegisterWatchEvent(erc20.ApprovalEvent(), erc20.TransferEvent())
	//e.RegisterWatchTopics(1, common.HexToHash("0xD7D761ce2e145FF4b72321F2075679ea42255286")) // filter from addr
	e.RegisterWatchTopics(2, common.HexToHash("0x9560d82d93a6a6d204df56f101964b26ce61e999")) // filter to addr
	e.RegisterEventHook(erc20.TransferEvent(), func(client *rpcclient.EvmClient, log types.Log) error {
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
	e.Close()
}

func TestQuickStartERC721(t *testing.T) {
	e, err := NewGeneralWatch(
		[]string{"https://sepolia.infura.io/v3/f6ef0da20fa14730ae77a316d88c0516"},
		[]common.Address{common.HexToAddress(ERC721)},
		&Options{
			Attrs: abs.Attrs{
				Chain:                "sepolia",
				DeployedBlockNumber:  6060135,
				ProcessedBlockNumber: 6152533,
				WatchBlockLimit:      50,
				ContractToDesc: map[string]abs.ContractDesc{
					ERC721: {
						Name:   "POPTAG",
						Symbol: "POP",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	e.RegisterWatchEvent(erc721.ApprovalEvent(), erc721.TransferEvent(), erc721.ApprovalForAllEvent())
	e.RegisterEventHook(erc721.TransferEvent(), func(client *rpcclient.EvmClient, log types.Log) error {
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
	e.Close()
}

func TestQuickStartTRON(t *testing.T) {
	evmAddr := utils.Tron.ToHexAddress("TEsYH363FySAj9UtDuc3ee6chU8DTLQqr3") // tvm convert to evm address
	e, err := NewGeneralWatch(
		[]string{"https://api.shasta.trongrid.io/jsonrpc"},
		[]common.Address{common.HexToAddress(evmAddr)}, // contract address
		&Options{
			Attrs: abs.Attrs{
				Chain:                "TRON",
				DeployedBlockNumber:  50701389,
				ProcessedBlockNumber: 50701389,
				WatchBlockLimit:      2,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	e.RegisterWatchEvent(abs.Event("0x4764e3effebe9df7fcf8e94a7f91735c90259220bfd64b99f636afae84cfc610"))
	e.RegisterEventHook(abs.Event("0x4764e3effebe9df7fcf8e94a7f91735c90259220bfd64b99f636afae84cfc610"), func(client *rpcclient.EvmClient, log types.Log) error {
		t.Logf("------ %d transfer log txhash: %s ------", log.BlockNumber, log.TxHash)
		t.Logf("topics: %v", log.Topics)
		t.Logf("data: %x", log.Data)
		return nil
	})

	for i := 0; i < 10; i++ {
		err = e.Watch()
		if err != nil {
			t.Fatal(err)
		}
	}
	e.Close()
}
