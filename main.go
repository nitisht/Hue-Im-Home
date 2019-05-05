package main

import (
	"fmt"
	"github.com/heatxsink/go-hue/configuration"
	"github.com/heatxsink/go-hue/lights"
	"github.com/heatxsink/go-hue/portal"
	"github.com/se1exin/hue-im-home/config"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// First check that the user has set the required env variables
	ipRange := os.Getenv("IP_RANGE")
	if ipRange == "" {
		fmt.Println("Required IP_RANGE environment variable not set.")
		fmt.Println("Please set IP_RANGE to an nmap style IP range")
		fmt.Println("E.g. 10.1.1.1,2,3 or 10.1.1.1-20")
		os.Exit(1)
	}

	appConfig := config.LoadConfig()
	if appConfig == nil {
		fmt.Println("Could not load config file! Exiting..")
		os.Exit(1)
	}

	// Check for a saved bridge, if not found look for a bridge on the network
	if appConfig.BridgeIPAddress == "" {
		fmt.Println("No bridge saved, looking for first bridge")

		pp, err := portal.GetPortal()
		if err != nil {
			fmt.Println("Failed to find any bridges on the network. Exiting..")
			os.Exit(1)
		}

		// We got a new bridge IP. Save it to the config for future
		appConfig.BridgeIPAddress = pp[0].InternalIPAddress
		config.SaveConfig(appConfig)

		fmt.Println("Found new bridge at " + appConfig.BridgeIPAddress)
	}

	// Check for a valid API key
	if appConfig.BridgeApiKey == "" {
		// API Key not found in config file, create a new API user on the bridge
		hueConfig := configuration.New(appConfig.BridgeIPAddress)
		response, err := hueConfig.CreateUser(appConfig.AppName, appConfig.AppName)
		if err != nil {
			fmt.Println("Error from bridge:", err)
			fmt.Println("Failed to register with your bridge.\nPlease ensure you have pressed the bridge button and try again.")
			os.Exit(1)
		}
		fmt.Println(response[0])
		// We got a new bridge Api Key. Save it to the config for future
		appConfig.BridgeApiKey = response[0].Success["username"].(string)
		config.SaveConfig(appConfig)
	}

	// All config steps passed. Now it's time to play with some lights..
	devicesOnline := scanForOnlineDevices()
	// Only attempt to change the lights if there has been a change since the last scan
	if appConfig.LastState != devicesOnline {
		fmt.Println("Device state changed to:", devicesOnline)
		if switchLights(appConfig, devicesOnline) {
			// We successfully changed the state, update the config for the next scan
			appConfig.LastState = devicesOnline
			config.SaveConfig(appConfig)
		}
	}
}

func scanForOnlineDevices() bool {
	// Run nmap, scanning port 5060 (androids) and 62078 (iphones) against the target ip range
	// Pipe the output of nmap into grep to find any occurances of 'open'.
	// If 'open' is found, then we have atleast one device online!
	// Anything else means that there are no devices online
	ipRange := os.Getenv("IP_RANGE")
	fmt.Println("Scanning network with range", ipRange)
	// out, err := exec.Command("nmap", "-p", "5060,62078", ipRange).Output()
	out, err := exec.Command("nmap", "-p", "62078", ipRange).Output()
	if err != nil {
		log.Fatal("Failed to run nmap!", err)
		return false
	}

	if strings.Contains(string(out), "open") {
		return true
	} else {
		return false
	}
}

func switchLights(appConfig *config.Config, newState bool) bool {
	hueLights := lights.New(appConfig.BridgeIPAddress, appConfig.BridgeApiKey)
	lightState := lights.State{On: newState}
	allLights, err := hueLights.GetAllLights()
	if err != nil {
		fmt.Println("Failed to get lights from Bridge!:", err)
		// Failed to connect to the bridge
		return false
	}

	// Turn on every single light
	// @TODO: Add a way to change only specific lights..
	for _, light := range allLights {
		// fmt.Printf("ID: %d Name: %s\n", light.ID, light.Name)
		/*
		if light.ID != 1 {
			continue // @TODO: Just testing for now..
		}
		*/
		_, err := hueLights.SetLightState(light.ID, lightState)
		if err != nil {
			fmt.Println("Failed to change light:", light.Name, err)
			return false
		}
	}

	return true
}