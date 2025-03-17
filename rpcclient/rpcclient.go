package rpcclient

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

type EvmClient struct {
	rawurl  string
	chainId uint64
	*ethclient.Client
}

func NewEvmRpcClient(rawurl string) (*EvmClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	id, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	return &EvmClient{rawurl: rawurl, chainId: id.Uint64(), Client: client}, nil
}

func MustNewEvmRpcClient(rawurl string) *EvmClient {
	c, err := NewEvmRpcClient(rawurl)
	if err != nil {
		err = fmt.Errorf("failed to dial %s, %v", rawurl, err)
		panic(err)
	}
	return c
}

func (c *EvmClient) GetRawUrl() string {
	return c.rawurl
}

func (c *EvmClient) GetChainId() uint64 {
	return c.chainId
}
