package main

import (
	"flag"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

var configFilePath = "./melodious.config.json"

func init() {
	flag.StringVar(&configFilePath, "config", configFilePath, "Path to the config file")
}

func main() {
	// Parse flags
	flag.Parse()

	// Load config
	cfg, err := NewConfigFromFile(configFilePath)
	if err != nil {
		panic(err)
	}

	// Setup logging
	log.SetHandler(text.New(os.Stderr))

	// Setup Melodious
	mel := NewMelodious(cfg)
	mel.ConnectToDB()
	<-mel.RunWebServer()
}
