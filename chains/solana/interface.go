package sol

import (
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/gagliardetto/solana-go"
)

type IContract interface {
	Init(attrs Attrs)
	Close() error
	DoneSignal() <-chan struct{}
	RegisterEventHook(event Event, f func(client *rpcclient.SolClient, txInfo TxInfo) error) error
	HandleEvent(client *rpcclient.SolClient, event Event, txInfo TxInfo) error
	UpdateProcessedTxSignature(txSig solana.Signature) error
	GetProcessedBlockNumber() solana.Signature
	Scan(client *rpcclient.SolClient) error
	GetContractDesc(programId string) (ContractDesc, error)
}
