package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config - stores
type Config struct {
	// Must be specified as shown at https://godoc.org/github.com/lib/pq
	DBAddr   string `json:"db-addr"`
	HTTPAddr string `json:"http-addr"`

	// These are ISO 8601 durations
	DeleteHistoryEvery string `json:"delete-history-every"`
	StoreHistoryFor    string `json:"store-history-for"`
}

// NewConfig - creates a new Config instance from given JSON data
func NewConfig(data []byte) (*Config, error) {
	var cfg Config
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// NewConfigFromFile - creates a new Config instance from given JSON file path
func NewConfigFromFile(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return NewConfig(data)
}
