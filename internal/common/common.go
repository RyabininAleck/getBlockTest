package common

import (
	"fmt"
	"math/big"
)

func IntToHex(i *big.Int) string {
	return fmt.Sprintf("0x%x", i)
}
