package sol

import (
	"fmt"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/AcSunday/gwatch-chain/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"testing"
)

type MintVoucherRecord struct {
	Ts             int64            `json:"ts"` // 时间戳
	Amount         uint64           `json:"amount"`
	TokenId        uint64           `json:"token_id"`
	MarketIndex    uint32           `json:"market_index"`
	SubAccountId   uint16           `json:"sub_account_id"`
	VoucherName    string           `json:"voucher_name"`
	Uri            string           `json:"uri"`
	User           solana.PublicKey `json:"user"`
	VoucherAddress solana.PublicKey `json:"voucher_address"`
}

func TestWatch(t *testing.T) {
	client, err := rpcclient.NewSolClient(rpc.DevNet_RPC, 1177777711)
	if err != nil {
		t.Fatal(err)
	}

	// https://explorer.solana.com/address/5ForPwE8sRNGy4vS5b7aWqHsAXbofAA8A89ZQ1QPq6Jk?cluster=devnet
	processedTxSignature := solana.MustSignatureFromBase58("4sDzE1GNKMSghfJh1HCUmv2voFbvWHVZ95GLbceTLDhzYyyBtEY4XxARhxcYiXpWpUvDoGpUcsfmnCh3hUqstYbX")
	contract, err := New("5ForPwE8sRNGy4vS5b7aWqHsAXbofAA8A89ZQ1QPq6Jk", &Attrs{
		ChainId:              client.GetChainId(),
		Chain:                "Solana",
		ProcessedTxSignature: processedTxSignature,
		WatchBlockLimit:      5,
	})
	if err != nil {
		t.Fatal(err)
	}

	//methodHash := utils.SigHashEvent(reflect.TypeOf(MintVoucherRecord{}).Name())
	methodHash := utils.GetMethodHashString(MintVoucherRecord{})
	err = contract.RegisterEventHook(Event(methodHash), func(client *rpcclient.SolClient, txInfo TxInfo) error {
		t.Logf("------ %d transfer event tx sig: %s ------", txInfo.TxDetail.Slot, txInfo.TxSig)
		t.Logf("tx info: %v", txInfo)

		var obj MintVoucherRecord
		err := utils.UnmarshalBorsh(txInfo.DataBytes, &obj)
		if err != nil {
			return fmt.Errorf("unmarshal tx info: %w", err)
		}
		t.Logf("obj: %+v", obj)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		err = contract.Watch(client)
		if err != nil {
			t.Fatal(err)
		}
	}
	contract.Close()
	client.Close()
}
