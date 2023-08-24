package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

type Configuration struct {
	Tz             string
	Token          string
	StartDate      time.Time
	UpdateInterval time.Duration
}

type Response struct {
	Data []DeviceData `json:"data"`
}

type DeviceData struct {
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

type DeviceInfo struct {
	TemperatureF int    `json:"temperatureF"`
	Humidity     int    `json:"humidity"`
	Ports        []Port `json:"ports"`
}

type Port struct {
	PortName string `json:"portName"`
	Speak    int    `json:"speak"`
	Port     int    `json:"port"`
}

func main() {
	log.Printf("Starting")

	file, err := os.Open("config.properties")
	if err != nil {
		log.Fatalf("Failed to open config file: %s", err)
	}
	defer file.Close()

	config := Configuration{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "=", 2) // Split at the first "="
		if len(split) != 2 {
			log.Fatalf("Invalid line in config: %s", line)
		}
		key, value := split[0], split[1]

		switch key {
		case "tz":
			config.Tz = value
		case "token":
			config.Token = value
		case "start_date":
			config.StartDate, err = time.Parse("2006-01-02", value)
			if err != nil {
				log.Fatalf("Failed to parse start_date: %s", err)
			}
		default:
			log.Printf("Unknown key: %s", key)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	for {
		location, err := time.LoadLocation(config.Tz)
		if err != nil {
			fmt.Println("Error loading location:", err)
			return
		}

		startDate := config.StartDate
		currentDate := time.Now().In(location)

		duration := currentDate.Sub(startDate)
		days := int(duration.Hours() / 24)

		formattedOutput := fmt.Sprintf("Day %d", days)
		err = ioutil.WriteFile("days.txt", []byte(formattedOutput), 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}

		url := "http://www.acinfinityserver.com/api/user/devInfoListAll?userId=" + config.Token
		reqBody := bytes.NewBuffer([]byte(""))

		req, err := http.NewRequest("POST", url, reqBody)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		req.Header.Add("token", config.Token)
		req.Header.Add("Host", "www.acinfinityserver.com")
		req.Header.Add("User-Agent", "okhttp/3.10.0")
		req.Header.Add("Content-Encoding", "gzip")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		var jsonResponse Response
		err = json.Unmarshal(respBody, &jsonResponse)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		if len(jsonResponse.Data) > 0 {
			temperatureF := float64(jsonResponse.Data[0].DeviceInfo.TemperatureF) / 100.0
			humidity := float64(jsonResponse.Data[0].DeviceInfo.Humidity) / 100.0

			formattedOutput := fmt.Sprintf("Temp: %.1f°", temperatureF)
			err = ioutil.WriteFile("temperatureF.txt", []byte(formattedOutput), 0644)
			if err != nil {
				log.Printf("Error writing to file: %v", err)
			}
			formattedOutput = fmt.Sprintf("Humidity: %.1f%%", humidity)
			err = ioutil.WriteFile("humidity.txt", []byte(formattedOutput), 0644)
			if err != nil {
				log.Printf("Error writing to file: %v", err)
			}
			log.Printf("Temperature: %.1fF, Humidity: %.1f%%", temperatureF, humidity)

			for _, deviceData := range jsonResponse.Data {
				for _, port := range deviceData.DeviceInfo.Ports {
					formattedOutput = fmt.Sprintf("%s: %d%%", port.PortName, port.Speak*10)
					err = ioutil.WriteFile("port_"+strconv.Itoa(port.Port)+".txt", []byte(formattedOutput), 0644)
					if err != nil {
						log.Printf("Error writing to file: %v", err)
					}
				}
			}
		}

		time.Sleep(config.UpdateInterval * time.Second)
	}
}
