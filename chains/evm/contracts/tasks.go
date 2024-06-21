package contracts

import "gwatch_chain/chains/evm/contracts/abs"

// TODO 通过tasks管理所有contract的扫描任务

type Tasks interface {
	// Register erc protocol watch task
	Register() error
}

type Options struct {
	WatchERC20   bool
	WatchERC721  bool
	WatchERC1155 bool
}

type tasks struct {
	abs.Contract
}

func NewTasks(opt *Options) Tasks {
	return &tasks{}
}

func (t *tasks) Register() error {
	return nil
}
