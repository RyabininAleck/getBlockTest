package server

import (
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"getBlockTest/internal/common"
	"getBlockTest/internal/models"
)

func (s *Server) analyzeBlocks(latestBlockNumber string) *models.AddressChange {
	latestBlockNumInt := s.parseBlockNumber(latestBlockNumber)
	balanceChanges := s.calculateBalanceChanges(latestBlockNumInt)
	return s.findMaxChange(balanceChanges)
}

func (s *Server) parseBlockNumber(latestBlockNumber string) *big.Int {
	latestBlockNum := new(big.Int)
	latestBlockNum.SetString(latestBlockNumber[2:], 16)
	return latestBlockNum
}

func (s *Server) calculateBalanceChanges(latestBlockNumInt *big.Int) map[string]*big.Int {
	balanceChanges := make(map[string]*big.Int)
	var blockNumInt *big.Int

	for i := 0; i < s.lengthLimit; i++ {
		blockNumInt = new(big.Int).Sub(latestBlockNumInt, big.NewInt(int64(i)))
		block, err := s.Adapter.GetBlockByNumber(common.IntToHex(blockNumInt))
		if err != nil {
			log.Println("Error getting block:", err)
			continue
		}

		s.processTransactions(block["transactions"].([]interface{}), balanceChanges)
	}
	return balanceChanges
}

func (s *Server) processTransactions(transactions []interface{}, balanceChanges map[string]*big.Int) {
	for _, tx := range transactions {
		txMap := tx.(map[string]interface{})

		from := txMap["from"].(string)
		to, toExists := txMap["to"].(string)
		valueHex := txMap["value"].(string)
		value := new(big.Int)
		value.SetString(valueHex[2:], 16)

		s.updateBalance(balanceChanges, from, value, false)
		if toExists && to != "" {
			s.updateBalance(balanceChanges, to, value, true)
		}
	}
}

func (s *Server) updateBalance(balanceChanges map[string]*big.Int, address string, value *big.Int, isAdd bool) {
	if balanceChanges[address] == nil {
		balanceChanges[address] = big.NewInt(0)
	}
	if isAdd {
		balanceChanges[address].Add(balanceChanges[address], value)
	} else {
		balanceChanges[address].Sub(balanceChanges[address], value)
	}
}

func (s *Server) findMaxChange(balanceChanges map[string]*big.Int) *models.AddressChange {
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

func (s *Server) GetMostChangedAddressHandler(w http.ResponseWriter, r *http.Request) {
	blockNumber, err := s.Adapter.GetLatestBlockNumber()
	if err != nil {
		log.Println("Error getting latest block number:", err)
		return
	}

	addressChange := s.analyzeBlocks(blockNumber)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(addressChange); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
