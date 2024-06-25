package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

func BigIntToDecimal(b *big.Int, decimalPlaces int32) decimal.Decimal {
	return decimal.NewFromBigInt(b, 0).Shift(-decimalPlaces)
}

func AmountToStr(_decimals int64, amount *big.Int) *big.Float {
	decimals := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(_decimals), nil)

	// 将代币总量转换为具有小数的十进制数
	tokenAmountFloat := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(decimals))

	return tokenAmountFloat
}

func EventAddressHashFormat(h common.Hash) string {
	return fmt.Sprintf("0x%s", h.String()[26:])
}
