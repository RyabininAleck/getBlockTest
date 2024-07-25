package adapter

import (
	"reflect"
	"testing"
)

func TestJSONRPCRequest_ToJSON(t *testing.T) {
	tests := []struct {
		name     string
		request  JSONRPCRequest
		expected string
		wantErr  bool
	}{
		{
			name: "valid request",
			request: JSONRPCRequest{
				JSONRPC: "2.0",
				Method:  "eth_blockNumber",
				Params:  []interface{}{},
				ID:      "getblock.io",
			},
			expected: `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":"getblock.io"}`,
			wantErr:  false,
		},
		{
			name: "request with params",
			request: JSONRPCRequest{
				JSONRPC: "2.0",
				Method:  "eth_getBlockByNumber",
				Params:  []interface{}{"latest", true},
				ID:      "getblock.io",
			},
			expected: `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest",true],"id":"getblock.io"}`,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &JSONRPCRequest{
				JSONRPC: tt.request.JSONRPC,
				Method:  tt.request.Method,
				Params:  tt.request.Params,
				ID:      tt.request.ID,
			}
			got, err := req.ToJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, []byte(tt.expected)) {
				t.Errorf("ToJSON() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
