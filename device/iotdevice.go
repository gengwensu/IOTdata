/*
iotdevice: a http client that simulates an iot device. It'll upload data bundles
to a HTTP server.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//IOTData - sensor data bundles; support only 2 types  "humidity" and "temperature"
type IOTData struct {
	DeviceUUID  string  `json:"uuid"`
	DeviceType  string  `json:"type"`
	SensorValue float64 `json:"sensor_value"`
	ReadTime    int64   `json:"sensor_reading_time"`
}

func main() {
	url := flag.String("url", "http://localhost:8081/iotData/add", "default url")
	uuid := flag.String("uuid", "device1", "device uuid")
	//deviceType := flag.String("type", "temperature", "device type")
	sValue := flag.Float64("sensorValue", 60.0, "60")
	tickTime := flag.Int("tick", 10, "Ticker duration in sec.")
	duration := flag.Int("duration", 300, "Total duration in sec.")
	//rTime := flag.Int64("readTime", time.Now().Unix(), "Unix sec.")
	flag.Parse()

	deviceType := "temperature"
	tickchan := time.NewTicker(time.Duration(*tickTime) * time.Second).C
	for {
		select {
		case <-time.After(time.Duration(*duration) * time.Second):
			fmt.Printf("device simulator %s is done\n", *uuid)
			break
		case <-tickchan:
			if deviceType == "temperature" { // alterating type for each tick
				deviceType = "humidity"
			} else {
				deviceType = "temperature"
			}
			rTime := time.Now().Unix()
			jsonString, err := json.Marshal(IOTData{*uuid, deviceType, *sValue, rTime})
			buf := strings.NewReader(string(jsonString))
			resp, err := http.Post(*url, "application/json", buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "device %s Post error: %v\n", *uuid, err)
				os.Exit(1)
			}
			_, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "device %s Post error: reading %s: %v\n", *uuid, *url, err)
				os.Exit(1)
			}
		}
	}
}
