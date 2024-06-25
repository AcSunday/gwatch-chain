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
