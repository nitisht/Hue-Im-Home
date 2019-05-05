package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const configPath = "/tmp/hue-im-home.json"

type Config struct {
	AppName string `json:"app_name"`
	BridgeApiKey string `json:"bridge_api_key"`
	BridgeIPAddress string `json:"bridge_ip_address"`
	LastState bool `json:"last_state"`
}

func LoadConfig() *Config {
	// Attempt to open the JSON file
	jsonFile, err := os.Open(configPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		// Could not open file, try to create it
		fmt.Println(err)
		return createNewConfig()
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// Try to parse the JSON. If it can't be parsed then we need to recreate it (corrupted file?)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	parseError := json.Unmarshal(byteValue, &config)
	if parseError != nil {
		fmt.Println(parseError)
		return createNewConfig()
	}

	return &config
}

func SaveConfig(config *Config) bool {
	file, _ := json.MarshalIndent(config, "", " ")
	err := ioutil.WriteFile(configPath, file, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func createNewConfig() *Config {
	fmt.Println("Creating new config file")
	config := Config {
		AppName: "Hue I'm Home",
		BridgeApiKey: "",
		BridgeIPAddress: "",
		LastState: false,
	}

	if SaveConfig(&config) {
		return &config
	}
	return nil
}