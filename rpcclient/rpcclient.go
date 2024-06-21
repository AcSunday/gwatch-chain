package rpcclient

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

type EvmClient struct {
	*ethclient.Client
}

func MustNewEvmRpcClient(rawurl string) *EvmClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rawurl)
	if err != nil {
		panic(err)
	}
	return &EvmClient{Client: client}
}
