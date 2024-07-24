package models

import (
	"encoding/json"
	"math/big"
)

type AddressChange struct {
	Address string   `json:"address"`
	Change  *big.Int `json:"change"`
}
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      string        `json:"id"`
}

func (req *JSONRPCRequest) ToJSON() ([]byte, error) {
	return json.Marshal(req)
}
