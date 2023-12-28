package main

import (
	"flag"
	"github.com/Seitanas/wosensortho-exporter/pkg/config"
	"github.com/Seitanas/wosensortho-exporter/pkg/prometheus"
	"log"
	"time"
)

func main() {
	var (
		configFile        string
		httpListenAddress string
		device            string
		btScanDuration    time.Duration
		btScanInterval    time.Duration
	)
	flag.StringVar(&configFile, "configFile", "/etc/wosensortho-exporter/config.json", "Path to config file.")
	flag.StringVar(&httpListenAddress, "httpListenAddress", "0.0.0.0:9353", "Address to bind to.")
	flag.StringVar(&device, "device", "default", "Implementation of ble")
	flag.DurationVar(&btScanDuration, "btScanDuration", time.Duration(5)*time.Second, "Duration in seconds for which exporter listens to sensor data")
	flag.DurationVar(&btScanInterval, "btScanInterval", time.Duration(15)*time.Second, "How often should exporter run sensor data listener")
	flag.Parse()
	err := config.Init(configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration. %v", err)
	}
	config.BTScanDuration = btScanDuration
	config.BTScanInterval = btScanInterval
	config.Device = device
	prometheus.Start(httpListenAddress)
}
