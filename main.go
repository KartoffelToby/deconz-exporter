package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	flag.StringVar(&token, "token", "", "The API token for deconz")
	flag.StringVar(&host, "host", "localhost", "The host address of the deconz instance")
	flag.IntVar(&port, "port", 0, "The port on which deconz is available")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")
}

func evalFlags() error {
	if token == "" {
		return fmt.Errorf("token is required")
	} else if host == "" {
		return fmt.Errorf("host is required")
	} else if port == 0 {
		return fmt.Errorf("port is required")
	}

	return nil
}

func main() {
	flag.Parse()
	if err := evalFlags(); err != nil {
		log.Fatalf("%s", err)
	}
	recordMetrics()
	serve()
}

func serve() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>Deconz Exporter</title></head>
            <body>
            <h1>Deconz Exporter</h1>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})
	fmt.Printf("Starting Server on %s:%d\n", "localhost", 2112)
	log.Fatal(http.ListenAndServe(":2112", nil))
}

var (
	host      = "192.168.0.222"
	port      = 1702
	token     = ""
	verbose   = false
	labels    = []string{"name", "uid", "manufacturer", "model", "type"}
	tmpMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "temperature",
		Help:      "Temperature of sensor in Celsius",
	}, labels)
	btryMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "battery",
		Help:      "Battery level of sensor in percent",
	}, labels)
	humidMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "humidity",
		Help:      "Humidity of sensor in percent",
	}, labels)
	pressureMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "pressure",
		Help:      "Air pressure in hectopascal (hPa)",
	}, labels)
	errorCtr = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "errors",
		Help:      "Failures to retrieve data from API",
	})
)

func recordMetrics() {
	url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   fmt.Sprintf("api/%s/sensors", token),
	}
	go func() {
		for {
			sensors, err := pollSensors(url)
			if err != nil {
				errorCtr.Inc()
				if verbose {
					fmt.Println("failed to retrieve data", err)
				}
			} else {
				if verbose {
					fmt.Printf("%+v", sensors)
				}
				for _, sensor := range sensors {
					labels := prometheus.Labels{"name": sensor.Name, "uid": sensor.UID, "manufacturer": sensor.Manufacturer, "model": sensor.ModelID, "type": sensor.Type}

					switch sensor.Type {
					case "ZHATemperature":
						btryMetric.With(labels).Set(float64(sensor.Config.Battery))
						tmpMetric.With(labels).Set(float64(sensor.State.Temperature) / 100.0)
					case "ZHAHumidity":
						humidMetric.With(labels).Set(float64(sensor.State.Humidity) / 100.0)
					case "ZHAPressure":
						pressureMetric.With(labels).Set(float64(sensor.State.Pressure))
					}
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func pollSensors(url url.URL) (map[string]Sensor, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var s map[string]Sensor
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, err
	}

	return s, err
}
