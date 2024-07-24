package server

import (
	"log"
	"net/http"

	"getBlockTest/internal/config"
)

type Adapter interface {
	GetLatestBlockNumber() (string, error)
	GetBlockByNumber(blockNumber string) (map[string]interface{}, error)
}

type Server struct {
	Port        string
	lengthLimit int
	Adapter     Adapter
}

func Create(config *config.Config, adapter Adapter) *Server {
	return &Server{Port: config.Port, lengthLimit: config.LengthLimit, Adapter: adapter}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/most-changed-address", s.GetMostChangedAddressHandler)

	handler := recoverMiddleware(loggingMiddleware(mux))

	log.Printf("Server is running on port %s\n", s.Port)
	err := http.ListenAndServe(s.Port, handler)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
