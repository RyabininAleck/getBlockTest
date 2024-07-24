package main

import (
	"getBlockTest/internal/adapter"
	"getBlockTest/internal/config"
	"getBlockTest/internal/server"
)

func main() {

	cfg := config.Load()
	adapt := adapter.Create(cfg)
	srv := server.Create(cfg, adapt)
	srv.Run()
}
