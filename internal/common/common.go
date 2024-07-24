package common

import (
	"fmt"
	"math/big"
)

func IntToHex(i *big.Int) string {
	if i.Sign() < 0 {
		return fmt.Sprintf("-0x%x", new(big.Int).Abs(i))
	}
	return fmt.Sprintf("0x%x", i)
}
