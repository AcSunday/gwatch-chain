package sol

import (
	"context"
	"errors"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"sync"
	"sync/atomic"
)

type Event string

const DefaultWatchLimit = 1000

type Attrs struct {
	ChainId      uint64
	Chain        string
	DeployedSlot uint64 // contract deployment slot

	// has been processed on tx signature
	//  can set the earliest transaction signature to start
	ProcessedTxSignature solana.Signature
	WatchBlockLimit      int                     // Limit the number of blocks scanned each time, default is 1000
	ContractToDesc       map[string]ContractDesc // key is programId
}

type ContractDesc struct {
	Name     string
	Symbol   string
	Decimals uint8
}

type TxInfo struct {
	ProgramId  solana.PublicKey
	TxSig      solana.Signature
	TxDetail   *rpc.GetTransactionResult
	DataBytes  []byte // ProgramData
	EventIndex uint64 // ProgramData index
}

type Contract struct {
	Attrs

	ProgramId solana.PublicKey

	IsRunning atomic.Bool
	IsClose   atomic.Bool

	handleFunc map[Event]func(client *rpcclient.SolClient, txInfo TxInfo) error // key is event
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func New(programId string, attrs *Attrs) (*Contract, error) {
	proId, err := solana.PublicKeyFromBase58(programId)
	if err != nil {
		return nil, err
	}

	if attrs.ProcessedTxSignature.IsZero() {
		return nil, errors.New("invalid processed tx signature, can set the earliest transaction signature to start")
	}

	c := &Contract{
		ProgramId: proId,
	}
	c.Init(*attrs)
	return c, nil
}

func (c *Contract) Init(attrs Attrs) {
	c.handleFunc = make(map[Event]func(client *rpcclient.SolClient, txInfo TxInfo) error, 4)

	c.Attrs = attrs
	if c.WatchBlockLimit <= 0 {
		c.WatchBlockLimit = DefaultWatchLimit
	}

	c.IsRunning.Store(true)
	c.IsClose.Store(false)
	c.mu = sync.RWMutex{}
	c.ctx, c.cancel = context.WithCancel(context.Background())
}

func (c *Contract) Close() error {
	if !c.IsClose.Load() {
		c.cancel()
		c.IsRunning.Store(false)
		c.IsClose.Store(true)
	}
	return nil
}

func (c *Contract) DoneSignal() <-chan struct{} {
	if c.IsRunning.Load() && !c.IsClose.Load() {
		return c.ctx.Done()
	}
	return nil
}

// RegisterEventHook Hook is a function that handles event,
// HandleEvent method call this Hook
func (c *Contract) RegisterEventHook(event Event, f func(client *rpcclient.SolClient, txInfo TxInfo) error) error {
	if c.IsClose.Load() {
		return errors.New("already closed, Registration of event hook is prohibited")
	}
	c.mu.Lock()
	c.handleFunc[event] = f
	c.mu.Unlock()
	return nil
}

// HandleEvent method call Hook
func (c *Contract) HandleEvent(client *rpcclient.SolClient, event Event, txInfo TxInfo) error {
	if !c.IsRunning.Load() {
		return errors.New("not running, handle event is prohibited")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if f, ok := c.handleFunc[event]; ok {
		return f(client, txInfo)
	}
	return nil
}

// UpdateProcessedTxSignature ...
func (c *Contract) UpdateProcessedTxSignature(txSig solana.Signature) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ProcessedTxSignature = txSig
	return nil
}

// GetProcessedTxSignature ...
func (c *Contract) GetProcessedTxSignature() solana.Signature {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ProcessedTxSignature
}

func (c *Contract) GetContractDesc(programId string) (ContractDesc, error) {
	v, ok := c.ContractToDesc[programId]
	if !ok {
		return ContractDesc{}, errors.New("not found")
	}
	return v, nil
}

func (e Event) String() string {
	return string(e)
}
