package utils

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"strings"
	"testing"
)

func TestDecimal(t *testing.T) {
	b, ok := big.NewInt(0).SetString("10000000000000000001", 10)
	if !ok {
		t.Fatal("failed to set integer")
	}
	d := BigIntToDecimal(b, 18)
	t.Log(d)
	t.Log(d.String())
}

func TestTronAddrConvert(t *testing.T) {
	addr := "35c63937d1165efa511f57838b050b9612e4b4a9"
	newAddr, err := Tron.FromHexAddress(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(newAddr)
	tronAddr := "TEsYH363FySAj9UtDuc3ee6chU8DTLQqr3"
	t.Log(Tron.ToHexAddress(tronAddr))
}

func TestParseEvmInputData(t *testing.T) {
	// 输入数据
	inputData := "0x806578880000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000083431333431323334000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a3534353332343532343500000000000000000000000000000000000000000000"
	// 将十六进制字符串转换为字节切片
	data, err := hex.DecodeString(inputData[10:])
	if err != nil {
		t.Fatal(err)
	}
	// 假设我们知道函数签名和参数类型，例如一个名为"getString"的函数，接受一个uint256参数
	methodABI, err := abi.JSON(strings.NewReader(`[{"name": "registerPID", "type": "function", "inputs": [{"name": "name_", "type": "string"},{"name": "code_", "type": "string"}], "outputs": []}]`))
	if err != nil {
		t.Fatal(err)
	}
	// 解析输入数据
	m := make(map[string]interface{})
	//t.Logf("data: %v", data)
	err = methodABI.Methods["registerPID"].Inputs.UnpackIntoMap(m, data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("map=== %v\n", m)
}
