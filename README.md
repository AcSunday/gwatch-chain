watch transaction events on the blockchain

go version 1.22.2

区块扫描，合约事件扫描

supports
  - evm
    - erc20
    - erc721

通过注册Hook function的方式，处理合约事件

支持topics过滤方式，过滤erc20类的from地址或to地址

简单用例请查看gwatch_test.go
