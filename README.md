# deconz-exporter

Exports measurements from Zigbee sensors registered with [deCONZ](https://dresden-elektronik.github.io/deconz-rest-doc/) as [Prometheus](https://prometheus.io/) metrics.

* Forked from [TobsCore/deconz-exporter](https://github.com/TobsCore/deconz-exporter).
* Supports temperature, humidity, pressure, door contact, thermostat and power consumption sensors.

## Configuration
An API token for the deCONZ REST API needs to be provided in the `DECONZ_TOKEN` environment variable.

Official instructions on this [here](https://dresden-elektronik.github.io/deconz-rest-doc/getting_started/#acquire-an-api-key). Having unlocked the gateway, the HTTP POST can be done easily with curl:
```
curl -d '{"devicetype":"deconz-exporter"}' -H "Content-Type: application/json" -X POST http://YOUR_DECONZ_HOST:YOUR_DECONZ_PORT/api
```

## Sample metrics
```
# HELP deconz_sensor_battery Battery level of sensor in percent
# TYPE deconz_sensor_battery gauge
deconz_sensor_battery{manufacturer="LUMI",model="lumi.weather",name="Aqara WSDCGQ11LM 0"} 100
# HELP deconz_sensor_errors Failures to retrieve data from API
# TYPE deconz_sensor_errors counter
deconz_sensor_errors 0
# HELP deconz_sensor_humidity Humidity of sensor in percent
# TYPE deconz_sensor_humidity gauge
deconz_sensor_humidity{manufacturer="LUMI",model="lumi.weather",name="Aqara WSDCGQ11LM 0",type="ZHAHumidity",uid="00:11:22:33:44:55:66:77-01-0405"} 58.63
# HELP deconz_sensor_pressure Air pressure in hectopascal (hPa)
# TYPE deconz_sensor_pressure gauge
deconz_sensor_pressure{manufacturer="LUMI",model="lumi.weather",name="Aqara WSDCGQ11LM 0",type="ZHAPressure",uid="00:11:22:33:44:55:66:77-01-0403"} 1015
# HELP deconz_sensor_sinceUpdate The time since the last update that was received from this sensor
# TYPE deconz_sensor_sinceUpdate gauge
deconz_sensor_sinceUpdate{manufacturer="LUMI",model="lumi.weather",name="Aqara WSDCGQ11LM 0"} 2672.645934682
# HELP deconz_sensor_temperature Temperature of sensor in Celsius
# TYPE deconz_sensor_temperature gauge
deconz_sensor_temperature{manufacturer="LUMI",model="lumi.weather",name="Aqara WSDCGQ11LM 0",type="ZHATemperature",uid="00:11:22:33:44:55:66:77-01-0402"} 21.08
...
```
