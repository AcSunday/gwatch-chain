package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
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
	//return fmt.Sprintf("0x%s", h.String()[26:])
	return common.HexToAddress(h.String()).String()
}

func ToValidateAddress(address string) string {
	addrLowerStr := strings.ToLower(address)
	if strings.HasPrefix(addrLowerStr, "0x") {
		addrLowerStr = addrLowerStr[2:]
		address = address[2:]
	}
	var binaryStr string
	addrBytes := []byte(addrLowerStr)
	hash256 := crypto.Keccak256Hash([]byte(addrLowerStr))
	for i, e := range addrLowerStr {
		if e >= '0' && e <= '9' {
			continue
		} else {
			binaryStr = fmt.Sprintf("%08b", hash256[i/2])
			if binaryStr[4*(i%2)] == '1' {
				addrBytes[i] -= 32
			}
		}
	}
	return "0x" + string(addrBytes)
}
