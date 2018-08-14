package lib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	messagebird "github.com/messagebird/go-rest-api"
)

type ServerConfig struct {
	ServerMode string
	Port       string
	AccessKey  string
}

var GlobalConfigFilePath = "./config.conf"
var GlobalConfig ServerConfig
var MBClient *messagebird.Client

func loadConfig() {
	configFile, err := os.Open(GlobalConfigFilePath)
	defer configFile.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(configFile)
	scanner.Split(bufio.ScanLines)

	var validTopic = regexp.MustCompile(`(\[)([a-zA-Z0-9]+)(\])`)
	var validContent = regexp.MustCompile(`([a-zA-Z0-9]+)(=)(([a-zA-Z0-9]+))`)

	currentPosition := ""

	for scanner.Scan() {

		if validTopic.MatchString(scanner.Text()) {
			topic := validTopic.ReplaceAllString(scanner.Text(), "$2")
			currentPosition = topic
		}

		if validContent.MatchString(scanner.Text()) {
			param := validContent.ReplaceAllString(scanner.Text(), "$1")

			if currentPosition == "server" {
				if param == "mode" {
					arg := validContent.ReplaceAllString(scanner.Text(), "$3")
					GlobalConfig.ServerMode = arg
				}

				if param == "port" {
					arg := validContent.ReplaceAllString(scanner.Text(), "$3")
					GlobalConfig.Port = arg
				}
			}

			if currentPosition == GlobalConfig.ServerMode {
				if param == "accessKey" {
					arg := validContent.ReplaceAllString(scanner.Text(), "$3")
					GlobalConfig.AccessKey = arg
				}
			}
		}
	}
}

func loadMessageBirdInstance() {
	MBClient = messagebird.NewV2(GlobalConfig.AccessKey)
}

func AppInitialize() {
	loadConfig()
	loadMessageBirdInstance()
}
