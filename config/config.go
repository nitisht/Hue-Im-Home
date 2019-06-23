package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	AppName string `json:"app_name"`
	BridgeApiKey string `json:"bridge_api_key"`
	BridgeIPAddress string `json:"bridge_ip_address"`
	LastState bool `json:"last_state"`
	LastOnTime string `json:"last_on_time"`
	LastOffTime string `json:"last_off_time"`
}

// Return the path to the config file
func configPath() string {
	path := "/config/config.json"

	if os.Getenv("CONFIG_FILE") != "" {
		path = os.Getenv("CONFIG_FILE")
	}

	return path
}

func LoadConfig() *Config {
	// Attempt to open the JSON file
	jsonFile, err := os.Open(configPath())
	// if we os.Open returns an error then handle it
	if err != nil {
		// Could not open file, try to create it
		log.Println(err)
		return createNewConfig()
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// Try to parse the JSON. If it can't be parsed then we need to recreate it (corrupted file?)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	parseError := json.Unmarshal(byteValue, &config)
	if parseError != nil {
		log.Println(parseError)
		return createNewConfig()
	}

	return &config
}

func SaveConfig(config *Config) bool {
	file, _ := json.MarshalIndent(config, "", " ")
	err := ioutil.WriteFile(configPath(), file, 0644)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func createNewConfig() *Config {
	log.Println("Creating new config file")
	config := Config {
		AppName: "Hue I'm Home",
		BridgeApiKey: "",
		BridgeIPAddress: "",
		LastState: false,
		LastOnTime: "",
		LastOffTime: "",
	}

	if SaveConfig(&config) {
		return &config
	}
	return nil
}

func (c *Config) GetLastOnTime() *time.Time {
	// Try to parse the Last On Time into a Time object
	_time, err := time.ParseInLocation("2006-01-02T15:04:05", c.LastOnTime, time.Local)
	if err != nil {
		return nil
	}
	return &_time
}

func (c *Config) GetLastOffTime() *time.Time {
	// Try to parse the Last Off Time into a Time object
	_time, err := time.ParseInLocation("2006-01-02T15:04:05", c.LastOffTime, time.Local)
	if err != nil {
		return nil
	}
	return &_time
}

func (c *Config) SetLastOnTime(time time.Time) {
	// Try to parse the Last On Time into a Time object
	c.LastOnTime = time.Format("2006-01-02T15:04:05")
}

func (c *Config) SetLastOffTime(time time.Time) {
	// Try to parse the Last On Time into a Time object
	c.LastOffTime = time.Format("2006-01-02T15:04:05")
}