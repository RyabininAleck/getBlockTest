package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Config struct {
	Port           string `json:"port"`
	ApiKey         string `json:"api_key"`
	LengthLimit    int    `json:"lengthLimit"`
	GetBlockParams struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      string `json:"id"`
	} `json:"getBlockParams"`
}

func Load() *Config {
	var filename string
	flag.StringVar(&filename, "file", "./config.json", "Path to the config file")
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	return &config
}
