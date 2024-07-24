package models

import "math/big"

type AddressChange struct {
	Address string   `json:"address"`
	Change  *big.Int `json:"change"`
}
