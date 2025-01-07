watch transaction events on the blockchain

go version 1.22.2

区块扫描，合约事件扫描

supports
  - evm
    - erc20
    - erc721
    - other...
    - 支持topics过滤方式，过滤erc20类的from地址或to地址
  - tvm
  - solana
    - 仅支持base64 decode
    - 仅支持新项目的扫描任务(官方rpc限制，推荐在已产生交易数量仅在几百条内使用)

通过注册Hook function的方式，处理合约事件

简单用例请查看gwatch_test.go
