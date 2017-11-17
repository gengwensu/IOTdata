/*
iotdata: a http server that accepts HTTP POST for uploading iot device data bundles to memory
and HTTP GET to display device data bundles already in memory.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
)

//IOTData - sensor data bundles; support only 2 types  "humidity" and "temperature"
type IOTData struct {
	DeviceUUID  string  `json:"uuid"`
	DeviceType  string  `json:"type"`
	SensorValue float64 `json:"sensor_value"`
	ReadTime    int64   `json:"sensor_reading_time"`
}

type dataStore []IOTData

const MAXOUTPUT = 50

var mux = &sync.Mutex{} //protect against race condition

func main() {
	//IOTdata in memory store
	ds := dataStore{}

	log.Fatal(http.ListenAndServe("localhost:8081", &ds))
}

func (db *dataStore) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/iotData/add": // POST
		if req.Method == "POST" {
			var inDataBundle IOTData
			err := json.NewDecoder(req.Body).Decode(&inDataBundle)
			if err != nil {
				http.Error(w, "Error decoding JSON request body.",
					http.StatusInternalServerError)
				//fmt.Fprintf(w, "inDataBundle %v\n", inDataBundle)
			}
			if inDataBundle.DeviceType == "humidity" ||
				inDataBundle.DeviceType == "temperature" { // a valid bundle
				mux.Lock()                      // guard against concurrent update
				*db = append(*db, inDataBundle) // Store in memory
				mux.Unlock()
			}

			fmt.Fprintf(w, "http %d\n", http.StatusOK)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	case "/iotData", "/iotData/": //GET
		if req.Method == "GET" {
			fmt.Fprint(w, "IOT data service\n") // return signature of the service
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	case "/iotData/deviceData": //GET deviceData
		if req.Method == "GET" {
			uuid := req.URL.Query().Get("uuid")
			deviceType := req.URL.Query().Get("type")
			startTime, err := strconv.ParseInt(req.URL.Query().Get("startTime"), 10, 64)
			if err != nil {
				startTime = 0
			}
			endTime, err := strconv.ParseInt(req.URL.Query().Get("endTime"), 10, 64)
			if err != nil {
				endTime = math.MaxInt64
			}

			out := []IOTData{}
			count := 0
			mux.Lock() // guard against concurrent update
			for _, r := range *db {
				if (r.DeviceUUID == uuid || uuid == "") &&
					(r.DeviceType == deviceType || deviceType == "") &&
					startTime <= r.ReadTime &&
					endTime >= r.ReadTime {
					out = append(out, r)
					count++
					if count >= MAXOUTPUT { //limit the number of output
						break
					}
				}
			}
			mux.Unlock()

			dataout, err := json.MarshalIndent(out, "", " ")
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}
			fmt.Fprintf(w, "results: %s\n", string(dataout))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "http %d, %s invalid. Only iotData, iotData/add, iotData/deviceData are allowed.\n",
			http.StatusMethodNotAllowed, req.URL)
	}
}
