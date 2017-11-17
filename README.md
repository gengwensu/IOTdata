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
    returns "IOT data service"

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
curl -X POST -H "application/json" --data "{"uuid": "b21ad0676f26439482cc9b1c7e827de4", "type": "temperature", "sensor_value": 60.0, "sensor_reading_time": 1510933249}" http://localhost:8081/IOTdata/add
{
  "ok"
}
```
The service should respond with 404 to all other requests not listed above

# environment & build
 require Go
  
$ go build ../src/github.com/gengwensu/IOTdata/iotdata.go 

$./iotData &

# iot device simulator

iotdevice.go in the device folder is a iot device simulator, which will send a HTTP POST to the server periodically for a duration. To run the simulator:

$ go build ../src/github.com/gengwensu/IOTdata/device/iotdevice.go

$ ./iotdevice -url="http://localhost:8081/iotData/add" ...

```
Usage of input flags:

  -duration int
        Total duration in sec. (default 300)
  -sensorValue float
        60 (default 60)
  -tick int
        Ticker duration in sec. (default 10)
  -url string
        default url (default "http://localhost:8081/iotData/add")
  -uuid string
        device uuid (default "device1")
```

# Running the experiment
```
$ ./iotdata.exe &
[1] 9224


$ ./iotdevice & ./iotdevice -uuid="device2" -tick=20 &
[2] 8260
[3] 9024

$ curl http://localhost:8081/iotData/deviceData
results: [
 {
  "uuid": "device1",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933179
 },
 {
  "uuid": "device1",
  "type": "temperature",
  "sensor_value": 60,
  "sensor_reading_time": 1510933189
 },
 {
  "uuid": "device2",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933189
 },
 {
  "uuid": "device1",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933199
 },
 {
  "uuid": "device2",
  "type": "temperature",
  "sensor_value": 60,
  "sensor_reading_time": 1510933209
 },
 {
  "uuid": "device1",
  "type": "temperature",
  "sensor_value": 60,
  "sensor_reading_time": 1510933209
 }
]

$ curl http://localhost:8081/iotData/deviceData?uuid=device1\&type=humidity\&startTime=1510933209\&endTime=1510933309
results: [
 {
  "uuid": "device1",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933219
 },
 {
  "uuid": "device1",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933239
 },
 {
  "uuid": "device1",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933259
 },
 {
  "uuid": "device1",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933279
 },
 {
  "uuid": "device1",
  "type": "humidity",
  "sensor_value": 60,
  "sensor_reading_time": 1510933299
 }
]

$ curl http://localhost:8081/iotData/deviceData?uuid=device2\&type=temperature\&startTime=1510933209\ &endTime=1510933309
results: [
 {
  "uuid": "device2",
  "type": "temperature",
  "sensor_value": 60,
  "sensor_reading_time": 1510933209
 },
 {
  "uuid": "device2",
  "type": "temperature",
  "sensor_value": 60,
  "sensor_reading_time": 1510933249
 },
 {
  "uuid": "device2",
  "type": "temperature",
  "sensor_value": 60,
  "sensor_reading_time": 1510933289
 }
]
```