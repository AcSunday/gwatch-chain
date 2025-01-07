package sol

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/AcSunday/gwatch-chain/utils"
	"github.com/gagliardetto/solana-go"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

var NotFoundProgramDataErr = errors.New("program data not found")

func (c *Contract) Watch(client *rpcclient.SolClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var txSigs []*rpc.TransactionSignature
	var err error

	txSigs, err = client.GetSignaturesForAddressWithOpts(ctx, c.ProgramId, &rpc.GetSignaturesForAddressOpts{
		//Limit:      &c.WatchBlockLimit,
		Until:      c.ProcessedTxSignature,
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return err
	}

	// handle tx sig
	for i := len(txSigs) - 1; i >= 0; i-- {
		err = c.handleTx(client, txSigs[i].Signature)
		if err != nil {
			return err
		}
	}

	// update ProcessedTxSignature
	if len(txSigs) == 0 {
		return nil
	}
	return c.UpdateProcessedTxSignature(txSigs[0].Signature)
}

func (c *Contract) handleTx(client *rpcclient.SolClient, txSig solana.Signature) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := client.GetTransaction(ctx, txSig, &rpc.GetTransactionOpts{
		Encoding:   solana.EncodingBase64,
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return fmt.Errorf("get transaction %v failed, %v", txSig, err)
	}
	if tx == nil || tx.Meta == nil {
		return fmt.Errorf("txSig: %v, get tx detail failed, tx or meta is nil", txSig)
	}

	programDatas, err := getDataBytesFromLogs(tx.Meta.LogMessages)
	if err != nil {
		return fmt.Errorf("txSig: %v, get data bytes from logs failed, %v", txSig, err)
	}

	for eventIdx, dataBytes := range programDatas {
		// Determine the name of the event
		if len(dataBytes) < 8 {
			continue
		}
		methodHash := dataBytes[:8]

		err := c.HandleEvent(client, Event(methodHash), TxInfo{
			ProgramId:  c.ProgramId,
			TxSig:      txSig,
			TxDetail:   tx,
			DataBytes:  dataBytes,
			EventIndex: uint64(eventIdx),
		})
		if err != nil {
			return err
		}

		//err = utils.UnmarshalBorsh(dataBytes, obj)
		//if err != nil {
		//	return err
		//}
	}

	return nil
}

func getDataBytesFromLogs(logs []string) ([][]byte, error) {
	programDatas := make([][]byte, 0, len(logs))
	for _, l := range logs {
		// Check whether the event starts with xx
		if !strings.HasPrefix(l, utils.ProgramDataPrefix) {
			continue
		}

		split := strings.Split(l, ":")
		if len(split) != 2 {
			continue
		}
		programData := strings.TrimSpace(split[1])

		dataBytes, err := base64.StdEncoding.DecodeString(programData)
		if err != nil {
			return nil, err
		}

		programDatas = append(programDatas, dataBytes)
	}

	if len(programDatas) == 0 {
		return nil, nil
	}
	return programDatas, nil
}
