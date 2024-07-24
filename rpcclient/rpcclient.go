package rpcclient

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rawurl)
	if err != nil {
		panic(err)
	}
	id, err := client.ChainID(ctx)
	if err != nil {
		panic(err)
	}
	return &EvmClient{rawurl: rawurl, chainId: id.Uint64(), Client: client}
}

func (c *EvmClient) GetRawUrl() string {
	return c.rawurl
}

func (c *EvmClient) GetChainId() uint64 {
	return c.chainId
}
