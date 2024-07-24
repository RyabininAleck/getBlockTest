package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"getBlockTest/internal/common"
	"getBlockTest/internal/models"
)

func (s *Server) GetMostChangedAddressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method not allowed. Method :", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	blockNumber, err := s.Adapter.GetLatestBlockNumber()
	if err != nil {
		log.Println("Error getting latest block number:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addressChange, err := s.analyzeBlocks(blockNumber)
	if err != nil {
		log.Println("Error analyzing block number:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(addressChange)
	if err != nil {
		log.Println("Error encoding address change:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) analyzeBlocks(latestBlockNumber string) (*models.AddressChange, error) {
	latestBlockNumInt := parseBlockNumber(latestBlockNumber)
	balanceChanges, err := s.calculateBalanceChanges(latestBlockNumInt)
	if err != nil {
		return nil, fmt.Errorf("error calculating balance changes: %w", err)
	}
	return findMaxChange(balanceChanges), nil
}

func (s *Server) calculateBalanceChanges(latestBlockNumInt *big.Int) (map[string]*big.Int, error) {
	balanceChanges := make(map[string]*big.Int)
	var blockNumInt *big.Int

	for i := 0; i < s.lengthLimit; i++ {
		blockNumInt = new(big.Int).Sub(latestBlockNumInt, big.NewInt(int64(i)))
		block, err := s.Adapter.GetBlockByNumber(common.IntToHex(blockNumInt))
		if err != nil {
			return nil, fmt.Errorf("error getting block:%w", err)
		}

		processTransactions(block["transactions"].([]interface{}), balanceChanges)
	}
	return balanceChanges, nil
}

func processTransactions(transactions []interface{}, balanceChanges map[string]*big.Int) {
	for _, tx := range transactions {
		txMap := tx.(map[string]interface{})

		from := txMap["from"].(string)
		to, toExists := txMap["to"].(string)
		valueHex := txMap["value"].(string)
		value := new(big.Int)
		value.SetString(valueHex[2:], 16)

		updateBalance(balanceChanges, from, value, false)
		if toExists && to != "" {
			updateBalance(balanceChanges, to, value, true)
		}
	}
}

func updateBalance(balanceChanges map[string]*big.Int, address string, value *big.Int, isAdd bool) {
	if balanceChanges[address] == nil {
		balanceChanges[address] = big.NewInt(0)
	}
	if isAdd {
		balanceChanges[address].Add(balanceChanges[address], value)
	} else {
		balanceChanges[address].Sub(balanceChanges[address], value)
	}
}

func parseBlockNumber(latestBlockNumber string) *big.Int {
	latestBlockNum := new(big.Int)
	latestBlockNum.SetString(latestBlockNumber[2:], 16)
	return latestBlockNum
}

func findMaxChange(balanceChanges map[string]*big.Int) *models.AddressChange {
	maxChangeAddress := ""
	maxChangeValue := new(big.Int)
	for addr, change := range balanceChanges {
		if change.Abs(change).Cmp(maxChangeValue.Abs(maxChangeValue)) > 0 {
			maxChangeAddress = addr
			maxChangeValue.Set(change)
		}
	}
	return &models.AddressChange{Address: maxChangeAddress, Change: maxChangeValue}
}
