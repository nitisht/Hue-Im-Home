package main

import (
	"github.com/heatxsink/go-hue/configuration"
	"github.com/heatxsink/go-hue/lights"
	"github.com/heatxsink/go-hue/portal"
	"github.com/se1exin/hue-im-home/config"
	"github.com/se1exin/hue-im-home/timerange"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	// First check that the user has set the required env variables
	ipRange := os.Getenv("IP_RANGE")
	if ipRange == "" {
		log.Println("Required IP_RANGE environment variable not set.")
		log.Println("Please set IP_RANGE to an nmap style IP range")
		log.Println("E.g. 10.1.1.1,2,3 or 10.1.1.1-20")
		os.Exit(1)
	}

	// Attempt to create/load the config file
	appConfig := config.LoadConfig()
	if appConfig == nil {
		log.Println("Could not create or load config file! Exiting..")
		os.Exit(1)
	}

	// Override the Bridge IP address if the user has set it via env
	bridgeIpAddress := os.Getenv("BRIDGE_IP")
	if bridgeIpAddress != "" {
		appConfig.BridgeIPAddress = bridgeIpAddress
	}

	// Start the run loop - which will run recursively over and over
	runLoop(appConfig)
}

func runLoop(appConfig *config.Config) {
	if checkConfig(appConfig) {

		onTimeRange := os.Getenv("ON_TIME_RANGE")
		offTimeRange := os.Getenv("OFF_TIME_RANGE")

		canSwitchOn := checkTimePeriodRipe(onTimeRange, appConfig.GetLastOnTime())
		canSwitchOff := checkTimePeriodRipe(offTimeRange, appConfig.GetLastOffTime())

		if canSwitchOn || canSwitchOff {
			// All config steps passed, and we are allowed to change the lights.
			// Check if devices are online or offline
			devicesOnline := scanForDevices()

			if devicesOnline && !canSwitchOn {
				log.Println("Device came online but we are not in valid ON_TIME_RANGE. Skipping..")
			} else if !devicesOnline && !canSwitchOff {
				log.Println("Device went offline but we are not in valid OFF_TIME_RANGE. Skipping..")
			} else if appConfig.LastState != devicesOnline {
				// Only attempt to change the lights if there has been a change since the last scan
				log.Println("Device state changed to:", devicesOnline)
				if switchLights(appConfig, devicesOnline) {
					// We successfully changed the state, update the config for the next scan
					appConfig.LastState = devicesOnline
					// Save the last on/off time for more advanced on/off logic later
					if devicesOnline {
						appConfig.SetLastOnTime(time.Now())
					} else {
						appConfig.SetLastOffTime(time.Now())
					}
					config.SaveConfig(appConfig)
				}
			}
		} else {
			log.Println("Cannot turn lights on or off. Nothing to do..")
		}
	}

	// Wait for SCAN_INTERVAL seconds, then run again!
	interval, err := strconv.Atoi(os.Getenv("SCAN_INTERVAL"))
	if err != nil {
		// Default to 10 seconds if the user did not pass a valid timeout period
		interval = 10
	}
	duration := time.Duration(interval) * time.Second

	time.Sleep(duration)

	runLoop(appConfig)
}

func checkConfig(appConfig *config.Config) bool {
	configValid := true

	// Check for a saved bridge, if not found look for a bridge on the network
	if appConfig.BridgeIPAddress == "" {
		log.Println("No bridge saved, looking for first bridge")

		pp, err := portal.GetPortal()
		if err != nil {
			log.Println("Failed to find any bridges on the network. Exiting..")
			configValid = false
		} else {
			// We got a new bridge IP. Save it to the config for future
			appConfig.BridgeIPAddress = pp[0].InternalIPAddress
			config.SaveConfig(appConfig)

			log.Println("Found new bridge at " + appConfig.BridgeIPAddress)
		}
	}

	// Check for a valid API key
	if appConfig.BridgeApiKey == "" {
		// API Key not found in config file, create a new API user on the bridge
		hueConfig := configuration.New(appConfig.BridgeIPAddress)
		response, err := hueConfig.CreateUser(appConfig.AppName, appConfig.AppName)
		if err != nil {
			log.Println("Error from bridge:", err)
			log.Println("Failed to register with your bridge.")
			log.Println("Please ensure you have pressed the bridge button.")
			configValid = false
		} else {
			// We got a new bridge Api Key. Save it to the config for future
			appConfig.BridgeApiKey = response[0].Success["username"].(string)
			log.Println("Successfully registered with bridge.")
			config.SaveConfig(appConfig)
		}
	}

	return configValid
}

/**
 * For a given Time Range (in format hh:mm-hh:mm), check that both the current time
 *  is within that range, and that a given comparison time is NOT within that range
 */
func checkTimePeriodRipe(inputTimeRange string, lastTime *time.Time) bool {
	timeRange, err := timerange.NewTimeRange(inputTimeRange)

	if err != nil {
		// User has not set an valid time range, default to always allow
		return true
	}

	// User has set an valid time range, but we are not currently in that range.
	if !timeRange.HasTime(time.Now()) {
		return false
	}

	// We are in the time range, have the lights been turned on/off in this time range before?
	if lastTime == nil {
		// The light has never been turned on/off.. so yes we can turn the lights on/off now
		return true
	}

	// Lights have been turned on/off before, but we cant turn them on/off again if already on/off in this time range
	return !timeRange.HasTime(*lastTime)
}

func scanForDevices() bool {
	devicesOnline := false
	// Run nmap, scanning port 5060 (androids) and 62078 (iphones) against the target ip range
	// Pipe the output of nmap into grep to find any occurances of 'open'.
	// If 'open' is found, then we have atleast one device online!
	// Anything else means that there are no devices online

	targetDevice := os.Getenv("DEVICE_TYPE") // Either 'android' or 'ios', or blank for both
	ipRange := os.Getenv("IP_RANGE")
	ports := "5060,62078" // default to both device types

	// If the user selected a specific device type, only scan those ports
	if targetDevice == "android" {
		ports = "5060"
	} else if targetDevice == "ios" {
		ports = "62078"
	} else {
		targetDevice = "all"
	}

	log.Printf("Scanning for %s devices on network in range %s", targetDevice, ipRange)

	start := time.Now()

	out, err := exec.Command("nmap", "-p", ports, ipRange).Output()
	if err != nil {
		log.Fatal("Failed to run nmap!", err)
	}

	if strings.Contains(string(out), "open") {
		devicesOnline = true
	}

	elapsed := time.Since(start)
	log.Printf("Scan completed in %s seconds", elapsed)

	return devicesOnline
}

func switchLights(appConfig *config.Config, newState bool) bool {
	hueLights := lights.New(appConfig.BridgeIPAddress, appConfig.BridgeApiKey)
	lightState := lights.State{On: newState}
	allLights, err := hueLights.GetAllLights()
	if err != nil {
		log.Println("Failed to get lights from Bridge!:", err)
		// Failed to connect to the bridge
		return false
	}

	// Turn on every single light
	// @TODO: Add a way to change only specific lights..
	for _, light := range allLights {
		// log.Printf("ID: %d Name: %s\n", light.ID, light.Name)
		/*
		if light.ID != 1 {
			continue // @TODO: Just testing for now..
		}
		*/
		_, err := hueLights.SetLightState(light.ID, lightState)
		if err != nil {
			log.Println("Failed to change light:", light.Name, err)
			return false
		}
	}

	return true
}