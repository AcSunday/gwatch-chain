package utils

import (
	"math/big"
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
