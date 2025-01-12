package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	serverURL        = "http://srv.msk01.gigacorp.local/_stats"
	maxLoadAverage   = 30
	memoryThreshold  = 0.8
	diskThreshold    = 0.9
	networkThreshold = 0.9
	checkInterval    = 10 * time.Second
	maxErrors        = 3
)

func main() {
	errorCount := 0

	for {
		resp, err := http.Get(serverURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			errorCount++
			if errorCount >= maxErrors {
				fmt.Println("Unable to fetch server statistic.")
				return
			}
			time.Sleep(checkInterval)
			continue
		}

		errorCount = 0 // Reset error count on successful response

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error reading response body:", err)
			time.Sleep(checkInterval)
			continue
		}

		data := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(data) != 7 {
			fmt.Println("Invalid data format received.")
			time.Sleep(checkInterval)
			continue
		}

		// Parse the data
		loadAverage, _ := strconv.ParseFloat(data[0], 64)
		totalMemory, _ := strconv.ParseFloat(data[1], 64)
		usedMemory, _ := strconv.ParseFloat(data[2], 64)
		totalDisk, _ := strconv.ParseFloat(data[3], 64)
		usedDisk, _ := strconv.ParseFloat(data[4], 64)
		totalNetwork, _ := strconv.ParseFloat(data[5], 64)
		usedNetwork, _ := strconv.ParseFloat(data[6], 64)

		// Check Load Average
		if loadAverage > maxLoadAverage {
			fmt.Printf("Load Average is too high: %.0f\n", loadAverage)
		}

		// Check Memory Usage
		memoryUsage := usedMemory / totalMemory
		if memoryUsage > memoryThreshold {
			fmt.Printf("Memory usage too high: %.0f%%\n", memoryUsage*100)
		}

		// Check Disk Space
		freeDisk := (totalDisk - usedDisk) / (1024 * 1024) // Convert to MB
		if usedDisk/totalDisk > diskThreshold {
			fmt.Printf("Free disk space is too low: %.0f Mb left\n", freeDisk)
		}

		// Check Network Bandwidth
		freeNetwork := usedNetwork / 8 / 1024 / 1024 // Convert to Mbit/s
		if usedNetwork/totalNetwork > networkThreshold {
			fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeNetwork)
		}

		time.Sleep(checkInterval)
	}
}
