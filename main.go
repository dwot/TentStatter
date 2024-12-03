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
	TemperatureF int      `json:"temperatureF"`
	Humidity     int      `json:"humidity"`
	Ports        []Port   `json:"ports"`
	Sensors      []Sensor `json:"sensors"`
}

type Port struct {
	PortName string `json:"portName"`
	Speak    int    `json:"speak"`
	Port     int    `json:"port"`
	CurMode  int    `json:"curMode"`
	Online   int    `json:"online"`
}

type Sensor struct {
	SensorType int `json:"sensorType"`
	AccessPort int `json:"accessPort"`
	SensorData int `json:"sensorData"`
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
			loc, _ := time.LoadLocation(config.Tz)
			config.StartDate, err = time.ParseInLocation("2006-01-02 15:04:05", value, loc)
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
		log.Printf("Current date: %v Start date: %v", currentDate, startDate)
		duration := currentDate.Sub(startDate)
		log.Printf("Duration: %v", duration)
		days := int(duration.Hours()/24) + 1
		log.Printf("Days: %d", days)
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
			for _, deviceData := range jsonResponse.Data {
				for _, sensor := range deviceData.DeviceInfo.Sensors {
					if sensor.SensorType == 0 { //Inside Temp
						value := float64(sensor.SensorData) / 100.0
						formattedOutput := fmt.Sprintf("Temp: %.1fÂ°", value)
						err = ioutil.WriteFile("inside_temp.txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}
					} else if sensor.SensorType == 2 { //Inside Humidity
						value := float64(sensor.SensorData) / 100.0
						formattedOutput = fmt.Sprintf("Humidity: %.1f%%", value)
						err = ioutil.WriteFile("inside_humidity.txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}
					} else if sensor.SensorType == 3 { //Inside VPD
						value := float64(sensor.SensorData) / 100.0
						formattedOutput = fmt.Sprintf("VPD: %.1f", value)
						err = ioutil.WriteFile("vpd.txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}
					} else if sensor.SensorType == 4 { //Outside Temp
						value := float64(sensor.SensorData) / 100.0
						formattedOutput := fmt.Sprintf("Temp: %.1fÂ°", value)
						err = ioutil.WriteFile("outside_temp.txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}
					} else if sensor.SensorType == 6 { //Outside Humidity
						value := float64(sensor.SensorData) / 100.0
						formattedOutput = fmt.Sprintf("Humidity: %.1f%%", value)
						err = ioutil.WriteFile("outside_humidity.txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}
						/*} else if sensor.SensorType == 7 { //Outside VPD
						value := float64(sensor.SensorData) / 100.0
						formattedOutput = fmt.Sprintf("VPD: %.1f", value)
						err = ioutil.WriteFile("outside_vpd.txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}*/
					}
				}

				for _, port := range deviceData.DeviceInfo.Ports {
					if port.Online == 1 {
						formattedOutput = fmt.Sprintf("%s: %d%%", port.PortName, port.Speak*10)
						err = ioutil.WriteFile("port_"+strconv.Itoa(port.Port)+".txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}
					} else {
						formattedOutput = "Offline"
						err = ioutil.WriteFile("port_"+strconv.Itoa(port.Port)+".txt", []byte(formattedOutput), 0644)
						if err != nil {
							log.Printf("Error writing to file: %v", err)
						}
					}
				}
			}

			log.Printf("Sleeping for 15 seconds", config.UpdateInterval)
			time.Sleep(15 * time.Second)
			log.Printf("Waking up")
		}
	}
}
