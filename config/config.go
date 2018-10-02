package config

import (
	"encoding/json"
	"os"

	messagebird "github.com/messagebird/go-rest-api"
)

var configFilePath = "./config.json"
var config map[string]map[string]string

// ServerMode is current server mode
var ServerMode string

// ServerPort is current server host port
var ServerPort string

// MBClient is MessageBird client instance
var MBClient *messagebird.Client

func loadConfig() {
	configFile, err := os.Open(configFilePath)

	if err != nil {
		panic("cannot open config file: " + err.Error())
	}

	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		panic("something went wrong when loading config file: " + err.Error())
	}

	ServerMode = config["server"]["mode"]
	ServerPort = config["server"]["port"]
}

func loadMessageBirdInstance() {
	MBClient = messagebird.New(config[ServerMode]["accessKey"])
}

// Initialize config
func Initialize() {
	loadConfig()
	loadMessageBirdInstance()
}
