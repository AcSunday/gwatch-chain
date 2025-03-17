package rpcclient

import (
	"fmt"
	"github.com/gagliardetto/solana-go/rpc"
)

type SolClient struct {
	rawurl  string
	chainId uint64
	*rpc.Client
}

// NewSolClient
// params: rawurl is rpc node
// params: chainId is diy chain id
func NewSolClient(rawurl string, chainId uint64) (*SolClient, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	client := rpc.New(rawurl)
	c := &SolClient{
		rawurl:  rawurl,
		chainId: chainId,
		Client:  client,
	}
	return c, nil
}

// MustNewSolClient
// params: rawurl is rpc node
// params: chainId is diy chain id
func MustNewSolClient(rawurl string, chainId uint64) *SolClient {
	c, err := NewSolClient(rawurl, chainId)
	if err != nil {
		err = fmt.Errorf("failed to dial %s, %v", rawurl, err)
		panic(err)
	}
	return c
}

func (c *SolClient) GetRawUrl() string {
	return c.rawurl
}

func (c *SolClient) GetChainId() uint64 {
	return c.chainId
}
