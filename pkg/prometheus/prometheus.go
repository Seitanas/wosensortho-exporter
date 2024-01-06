package prometheus

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"sync"
	"time"
	//    "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/Seitanas/wosensortho-exporter/pkg/btle"
	"github.com/Seitanas/wosensortho-exporter/pkg/config"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mutex = sync.RWMutex{}
)

func btScan() {
	d, err := dev.NewDevice(config.Device)
	if err != nil {
		log.Fatalf("Can't create device: %s", err)
	}
	ble.SetDefaultDevice(d)
	defer ble.Stop()
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), config.BTScanDuration))
	chkErr(ble.Scan(ctx, true, btle.Handler, nil))
}

func recordMetrics() {
	go func() {
		for {
			btScan()
			time.Sleep(config.BTScanInterval)
		}
	}()
}

func staticPage(w http.ResponseWriter, req *http.Request) {
	page := `<html>
    <head><title>wosensortho exporter</title></head>
    <body>
    <h1>wosensortho exporter</h1>
    <p><a href='metrics'>Metrics</a></p>
    </body>
    </html>`
	fmt.Fprintln(w, page)
}

type sensorCollector struct {
}

func newSensorCollector() *sensorCollector {
	return &sensorCollector{}
}

func (collector *sensorCollector) Describe(ch chan<- *prometheus.Desc) {

}

func buildPromDesc(name string, description string, labels map[string]string) *prometheus.Desc {
	return prometheus.NewDesc(
		name,
		description,
		nil,
		labels,
	)
}

func (collector *sensorCollector) Collect(ch chan<- prometheus.Metric) {

	var desc *prometheus.Desc
	mutex.Lock()
	defer mutex.Unlock()

	for mac, data := range btle.BTDevice {
		labels := make(map[string]string)
		labels["mac"] = mac
		for _, l := range config.Config.Sensors[mac].Labels {
			labels[l.Name] = l.Value
		}
		if len(data.ManufacturerData) > 0 {
			desc = buildPromDesc("wosensortho_temperature", "Temperature reading", labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, data.Temperature)

			desc = buildPromDesc("wosensortho_humidity", "Humidity reading", labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, data.Humidity)
		}
		if len(data.ServiceData) > 0 {
			desc = buildPromDesc("wosensortho_battery", "Battery level", labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, data.Battery)
		}
	}
}

func Start(httpListenAddress string) {
	recordMetrics()
	sensor := newSensorCollector()
	prometheus.MustRegister(sensor)
	router := mux.NewRouter()
	router.HandleFunc("/", staticPage)
	http.Handle("/", router)
	router.Path("/metrics").Handler(promhttp.Handler())
	err := http.ListenAndServe(httpListenAddress, router)
	log.Fatal(err)
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
	case context.Canceled:
		log.Printf("Canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
