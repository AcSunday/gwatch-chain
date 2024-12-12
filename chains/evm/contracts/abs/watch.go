package abs

import (
	"context"
	"github.com/AcSunday/gwatch-chain/rpcclient"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"time"
)

func (c *Contract) Watch(client *rpcclient.EvmClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get latest block
	latestNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return err
	}
	if c.ProcessedBlockNumber+1 > latestNumber {
		return nil
	}

	startBlockNumber := int64(c.ProcessedBlockNumber + 1)
	endBlockNumber := startBlockNumber + c.WatchBlockLimit
	if endBlockNumber > int64(latestNumber) {
		endBlockNumber = int64(latestNumber)
	}

	// filter data on the chain
	query := c.getFilterQuery(startBlockNumber, endBlockNumber)
	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	for _, l := range logs {
		// filter not have topic, or has been reverted
		if len(l.Topics) == 0 || l.Removed {
			continue
		}

		event := l.Topics[0]
		err := c.HandleEvent(client, Event(event.Hex()), l)
		if err != nil {
			return err
		}
	}

	c.UpdateProcessedBlockNumber(uint64(endBlockNumber))

	return nil
}

func (c *Contract) getFilterQuery(startBlockNumber, endBlockNumber int64) ethereum.FilterQuery {
	query := ethereum.FilterQuery{
		Addresses: c.Addrs,
		FromBlock: big.NewInt(startBlockNumber),
		ToBlock:   big.NewInt(endBlockNumber),
	}

	if len(c.Topics) > 0 {
		c.mu.RLock()
		query.Topics = c.Topics
		c.mu.RUnlock()
		//query.Topics = make([][]common.Hash, 4)
		//query.Topics[0] = append(c.Topics, topics[0]...)
		//query.Topics[2] = []common.Hash{common.HexToHash("0x59330ab2485985a1cd76cb0239bd37378978b0ea"),
		//	common.HexToHash("0x73bf6617837d6ada5e5a48c1017af46b016e2dcb")}
	}

	return query
}
