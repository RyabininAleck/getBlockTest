package server

import (
	"math/big"
	"testing"
)

func TestFindMaxChange(t *testing.T) {
	tests := []struct {
		name           string
		balanceChanges map[string]*big.Int
		expected       *AddressChange
	}{
		{
			name: "single address",
			balanceChanges: map[string]*big.Int{
				"0x1": big.NewInt(100),
			},
			expected: &AddressChange{
				Address: "0x1",
				Change:  big.NewInt(100),
			},
		},
		{
			name: "multiple addresses, positive changes",
			balanceChanges: map[string]*big.Int{
				"0x1": big.NewInt(100),
				"0x2": big.NewInt(200),
				"0x3": big.NewInt(150),
			},
			expected: &AddressChange{
				Address: "0x2",
				Change:  big.NewInt(200),
			},
		},
		{
			name: "multiple addresses, mixed changes",
			balanceChanges: map[string]*big.Int{
				"0x1": big.NewInt(-100),
				"0x2": big.NewInt(200),
				"0x3": big.NewInt(-300),
				"0x4": big.NewInt(150),
			},
			expected: &AddressChange{
				Address: "0x3",
				Change:  big.NewInt(300),
			},
		},
		{
			name: "multiple addresses, same changes",
			balanceChanges: map[string]*big.Int{
				"0x1": big.NewInt(100),
				"0x2": big.NewInt(100),
				"0x3": big.NewInt(100),
			},
			expected: &AddressChange{
				Address: "0x1",
				Change:  big.NewInt(100),
			},
		},
		{
			name:           "empty map",
			balanceChanges: map[string]*big.Int{},
			expected: &AddressChange{
				Address: "",
				Change:  big.NewInt(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findMaxChange(tt.balanceChanges)
			if result.Address != tt.expected.Address {
				t.Errorf("findMaxChange() Address = %v, want %v", result.Address, tt.expected.Address)
			}
			if result.Change.Cmp(tt.expected.Change) != 0 {
				t.Errorf("findMaxChange() Change = %v, want %v", result.Change, tt.expected.Change)
			}
		})
	}
}

func TestParseBlockNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *big.Int
	}{
		{
			name:     "zero",
			input:    "0x0",
			expected: big.NewInt(0),
		},
		{
			name:     "small number",
			input:    "0xf",
			expected: big.NewInt(15),
		},
		{
			name:     "large number",
			input:    "0x75bcd15",
			expected: big.NewInt(123456789),
		},
		{
			name:  "very large number",
			input: "0x123456789abcdef",
			expected: func() *big.Int {
				val, _ := new(big.Int).SetString("123456789abcdef", 16)
				return val
			}(),
		},
		{
			name:     "invalid input",
			input:    "0xGHIJK",
			expected: big.NewInt(0),
		},
		{
			name:     "empty input",
			input:    "0x",
			expected: big.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBlockNumber(tt.input)
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("parseBlockNumber(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUpdateBalance(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]*big.Int
		address  string
		value    *big.Int
		isAdd    bool
		expected map[string]*big.Int
	}{
		{
			name: "add balance to existing address",
			initial: map[string]*big.Int{
				"0x1": big.NewInt(100),
			},
			address: "0x1",
			value:   big.NewInt(50),
			isAdd:   true,
			expected: map[string]*big.Int{
				"0x1": big.NewInt(150),
			},
		},
		{
			name: "subtract balance from existing address",
			initial: map[string]*big.Int{
				"0x1": big.NewInt(100),
			},
			address: "0x1",
			value:   big.NewInt(50),
			isAdd:   false,
			expected: map[string]*big.Int{
				"0x1": big.NewInt(50),
			},
		},
		{
			name:    "add balance to new address",
			initial: map[string]*big.Int{},
			address: "0x2",
			value:   big.NewInt(75),
			isAdd:   true,
			expected: map[string]*big.Int{
				"0x2": big.NewInt(75),
			},
		},
		{
			name:    "subtract balance from new address",
			initial: map[string]*big.Int{},
			address: "0x3",
			value:   big.NewInt(30),
			isAdd:   false,
			expected: map[string]*big.Int{
				"0x3": big.NewInt(-30),
			},
		},
		{
			name: "add zero balance",
			initial: map[string]*big.Int{
				"0x1": big.NewInt(100),
			},
			address: "0x1",
			value:   big.NewInt(0),
			isAdd:   true,
			expected: map[string]*big.Int{
				"0x1": big.NewInt(100),
			},
		},
		{
			name: "subtract zero balance",
			initial: map[string]*big.Int{
				"0x1": big.NewInt(100),
			},
			address: "0x1",
			value:   big.NewInt(0),
			isAdd:   false,
			expected: map[string]*big.Int{
				"0x1": big.NewInt(100),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Копируем начальные данные, чтобы избежать изменений в оригинале
			balanceChanges := make(map[string]*big.Int)
			for k, v := range tt.initial {
				balanceChanges[k] = new(big.Int).Set(v)
			}

			updateBalance(balanceChanges, tt.address, tt.value, tt.isAdd)

			for addr, expectedValue := range tt.expected {
				if resultValue, exists := balanceChanges[addr]; !exists {
					t.Errorf("address %v not found in balanceChanges", addr)
				} else if resultValue.Cmp(expectedValue) != 0 {
					t.Errorf("balance for address %v = %v, want %v", addr, resultValue, expectedValue)
				}
			}
		})
	}
}

func TestProcessTransactions(t *testing.T) {
	tests := []struct {
		name         string
		transactions []interface{}
		expected     map[string]*big.Int
	}{
		{
			name: "single transaction",
			transactions: []interface{}{
				map[string]interface{}{
					"from":  "0x1",
					"to":    "0x2",
					"value": "0xa",
				},
			},
			expected: map[string]*big.Int{
				"0x1": big.NewInt(-10),
				"0x2": big.NewInt(10),
			},
		},
		{
			name: "multiple transactions",
			transactions: []interface{}{
				map[string]interface{}{
					"from":  "0x1",
					"to":    "0x2",
					"value": "0xa",
				},
				map[string]interface{}{
					"from":  "0x2",
					"to":    "0x3",
					"value": "0x14",
				},
			},
			expected: map[string]*big.Int{
				"0x1": big.NewInt(-10),
				"0x2": big.NewInt(-10),
				"0x3": big.NewInt(20),
			},
		},
		{
			name: "transaction without recipient",
			transactions: []interface{}{
				map[string]interface{}{
					"from":  "0x1",
					"value": "0xa",
				},
			},
			expected: map[string]*big.Int{
				"0x1": big.NewInt(-10),
			},
		},
		{
			name:         "empty transactions",
			transactions: []interface{}{},
			expected:     map[string]*big.Int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balanceChanges := make(map[string]*big.Int)
			processTransactions(tt.transactions, balanceChanges)

			for addr, expectedValue := range tt.expected {
				if resultValue, exists := balanceChanges[addr]; !exists {
					t.Errorf("address %v not found in balanceChanges", addr)
				} else if resultValue.Cmp(expectedValue) != 0 {
					t.Errorf("balance for address %v = %v, want %v", addr, resultValue, expectedValue)
				}
			}
		})
	}
}
