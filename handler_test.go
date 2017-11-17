package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var bundle = []IOTData{
	{"device1", "temperature", 50.0, 1510933479},
	{"device2", "humidity", 80.0, 1510933489},
}

func TestPOSTSaveBundle(t *testing.T) { // POST a bundle
	var tests = []struct {
		payload IOTData
		want    dataStore
	}{
		{ //test1
			bundle[0],
			dataStore{bundle[0]},
		},
		{ //test2
			bundle[1],
			dataStore{bundle[0], bundle[1]},
		},
	}

	handler := &dataStore{}
	server := httptest.NewServer(handler)
	defer server.Close()

	for _, test := range tests {
		jsonString, err := json.Marshal(test.payload)
		buf := strings.NewReader(string(jsonString))
		url := server.URL + "/iotData/add"
		resp, err := http.Post(url, "application/json", buf)
		if err != nil {
			t.Fatal(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("HTTP returns status code %d; want %d\n", resp.StatusCode, http.StatusOK)
		}

		if !reflect.DeepEqual(*handler, test.want) {
			t.Errorf("payload %v not saved properly, want %v\n", *handler, test.want)
		}
	}
}

func TestGETListBundles(t *testing.T) { // GET - list bundles
	b1, _ := json.MarshalIndent([]IOTData{bundle[0]}, "", " ")
	b2, _ := json.MarshalIndent([]IOTData{bundle[1]}, "", " ")
	b3, _ := json.MarshalIndent([]IOTData{}, "", " ")
	s1 := fmt.Sprintf("results: %s\n", string(b1))
	s2 := fmt.Sprintf("results: %s\n", string(b2))
	s3 := fmt.Sprintf("results: %s\n", string(b3))

	var tests = []struct {
		params map[string]string
		want   string
	}{
		{
			map[string]string{
				"uuid":      "device1",
				"type":      "temperature",
				"startTime": "1510933479",
				"endTime":   "1510933489",
			},
			s1,
		},
		{
			map[string]string{
				"uuid":      "device2",
				"type":      "humidity",
				"startTime": "1510933479",
				"endTime":   "1510933489",
			},
			s2,
		},
		{
			map[string]string{
				"uuid":      "device2",
				"type":      "humidity",
				"startTime": "1510933469",
				"endTime":   "1510933479",
			},
			s3,
		},
		{
			map[string]string{
				"uuid":      "device2",
				"type":      "humidity",
				"startTime": "1510933499",
				"endTime":   "1510933509",
			},
			s3,
		},
	}

	handler := &dataStore{bundle[0], bundle[1]}
	server := httptest.NewServer(handler)
	defer server.Close()

	for _, test := range tests {
		baseURL := server.URL + "/iotData/deviceData?"
		p := url.Values{}
		for k, v := range test.params {
			p.Add(k, v)
		}
		url := baseURL + p.Encode()
		resp, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		result, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("HTTP returns status code %d; want %d\n", resp.StatusCode, http.StatusOK)
		}

		if string(result) != test.want {
			t.Errorf("result %v incorrect, want %v\n", string(result), test.want)
		}
	}
}
