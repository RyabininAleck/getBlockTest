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
	getLastBlockBody        = `{"jsonrpc": "2.0","method": "eth_blockNumber","params": [],"id": "getblock.io"}`
	getBlockByNumBodyFormat = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":"getblock.io"}`
)

func Create(config *config.Config) *Adapter {
	return &Adapter{apiKey: config.APIKey}
}

type Adapter struct {
	apiKey string
}

func (a *Adapter) GetLatestBlockNumber() (string, error) {
	getBlockNumberURL := fmt.Sprintf(getBlockNumberURLFormat, a.apiKey)
	req, err := http.NewRequest("POST", getBlockNumberURL, bytes.NewBuffer([]byte(getLastBlockBody)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result["result"].(string), nil
}

func (a *Adapter) GetBlockByNumber(blockNumber string) (map[string]interface{}, error) {
	getBlockNumberURL := fmt.Sprintf(getBlockNumberURLFormat, a.apiKey)
	payload := []byte(fmt.Sprintf(getBlockByNumBodyFormat, blockNumber))

	req, err := http.NewRequest("POST", getBlockNumberURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result["result"].(map[string]interface{}), nil
}
