# IOTdata

Simple Data service for IOT devices, no auth. no db; memory only

# Usecase

* IOT devices can concurrently upload to this service with sensor data bundles, a JSON string that looks like this:


```
{
"device_uuid": "b21ad0676f26439482cc9b1c7e827de4",
"sensor_type": "temperature",
"sensor_value": 50,
"sensor_reading_time": 1510093202
}
```

There are 2 valid sensor_types: "humidity" and "temperature".

* Clients can make requests to this service to retrieve sensor data for a given device in a time range with the following JSON 
bundle: 

```
{
"device_uuid": "b21ad0676f26439482cc9b1c7e827de4",
"sensor_type": "temperature",
"start_time": 1510093202,
"end_time" 1510099302
}
```

# API

The IOTdata service will run on http://localhost:8081 and will support the following RESTful APIs:

1. GET /iotData/
    returns "url IOT data service"

2. GET /iotData/deviceData?uuid="b21ad0676f26439482cc9b1c7e827de4"&type="temperature"&startTime=1510093202&endTime=1510099202
    returns a JSON string of bundles

    example: 


    ```
    $ curl http://localhost:8081/iotData/deviceData?uuid="b21ad0676f26439482cc9b1c7e827de4"&type="temperature"&startTime=1510093202&endTime=1510099202
    {
      {
        "device_uuid": "b21ad0676f26439482cc9b1c7e827de4",
        "sensor_type": "temperature",
        "sensor_value": 50,
        "sensor_reading_time": 1510095202
      }
    }
    
    ```
 
3. POST /add

A request to add a bundle. When the request is processed by the server, the bundle will be added into memory.

Example:

```
curl -X POST --data "$device_bundle_data_in_JSON_format" http://localhost:8081/iotData/add
{
  "ok"
}
```
The service should respond with 404 to all other requests not listed above

# environment & build
 require Go
  
$ go build ../src/github.com/gengwensu/iotData/iotdata.go 

$./iotData &
...

