package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Config struct {
	Port        string `json:"port"`
	LengthLimit int    `json:"lengthLimit"`
	APIKey      string `json:"api_key"`
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
