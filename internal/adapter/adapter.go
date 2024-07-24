package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"getBlockTest/internal/config"
)

const (
	getBlockNumberURLFormat = "https://go.getblock.io/%s/"
	getLastBlockMethod      = "eth_blockNumber"
	getBlockByNumMethod     = "eth_getBlockByNumber"
)

func Create(config *config.Config) *Adapter {
	return &Adapter{apiKey: config.ApiKey, jsonRPCVersion: config.GetBlockParams.Jsonrpc, id: config.GetBlockParams.Id}
}

type Adapter struct {
	apiKey         string
	jsonRPCVersion string
	id             string
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

func (a *Adapter) GetBlockByNumber(blockNumber string) (map[string]interface{}, error) {
	requestBody := &JSONRPCRequest{
		JSONRPC: a.jsonRPCVersion,
		Method:  getBlockByNumMethod,
		Params:  []interface{}{blockNumber, true},
		ID:      a.id,
	}

	result, err := a.sendRequest(requestBody)
	if err != nil {
		return nil, err
	}

	blockResult, ok := result["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: %v", result)
	}

	return blockResult, nil
}

func (a *Adapter) GetLatestBlockNumber() (string, error) {
	requestBody := &JSONRPCRequest{
		JSONRPC: a.jsonRPCVersion,
		Method:  getLastBlockMethod,
		Params:  []interface{}{},
		ID:      a.id,
	}

	result, err := a.sendRequest(requestBody)
	if err != nil {
		return "", err
	}

	blockNumber, ok := result["result"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format: %v", result)
	}

	return blockNumber, nil
}

func (a *Adapter) sendRequest(body *JSONRPCRequest) (map[string]interface{}, error) {
	jsonData, err := body.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("Error creating JSON: %w", err)
	}

	url := fmt.Sprintf(getBlockNumberURLFormat, a.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from server", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
