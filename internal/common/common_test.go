package common

import (
	"math/big"
	"testing"
)

func TestIntToHex(t *testing.T) {

	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{
			name:     "zero",
			input:    big.NewInt(0),
			expected: "0x0",
		},
		{
			name:     "small number",
			input:    big.NewInt(15),
			expected: "0xf",
		},
		{
			name:     "large number",
			input:    big.NewInt(123456789),
			expected: "0x75bcd15",
		},
		{
			name: "very large number",
			input: func() *big.Int {
				val, _ := new(big.Int).SetString("123456789abcdef", 16)
				return val
			}(),
			expected: "0x123456789abcdef",
		},
		{
			name:     "negative number",
			input:    big.NewInt(-123456789),
			expected: "-0x75bcd15",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToHex(tt.input); got != tt.expected {
				t.Errorf("IntToHex() = %v, want %v", got, tt.expected)
			}
		})
	}
}
