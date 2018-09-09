package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	messagebird "github.com/messagebird/go-rest-api"
)

const globalConfigFilePath = "./config.json"

var (
	config     map[string]map[string]string
	serverMode string
	serverPort string
	mbClient   *messagebird.Client
)

func loadConfig() {
	configFile, err := os.Open(globalConfigFilePath)
	defer configFile.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Println("Error loading config file", err)
	}

	serverMode = config["server"]["mode"]
	serverPort = config["server"]["port"]
}

func loadMessageBirdInstance() {
	mbClient = messagebird.New(config[serverMode]["accessKey"])
}

func initConfig() {
	loadConfig()
	loadMessageBirdInstance()
}
