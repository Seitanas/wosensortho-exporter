package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/Seitanas/wosensortho-exporter/pkg/btle"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"
	"log"
	"time"
)

var (
	device = flag.String("device", "default", "Implementation of ble")
	du     = flag.Duration("du", 5*time.Second, "Scan duration")
)

func main() {
	flag.Parse()

	d, err := dev.NewDevice(*device)
	if err != nil {
		log.Fatalf("Can't create device: %s", err)
	}
	ble.SetDefaultDevice(d)

	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, true, btle.Handler, nil))

	fmt.Printf("Found SwitchBot devices:\n")
	for mac, data := range btle.BTDevice {
		fmt.Printf("MAC: %s, ManufacturerData: %s ServiceData: %s Temperature: %f Humidity: %f Battery: %f\n", mac, hex.EncodeToString(data.ManufacturerData), hex.EncodeToString(data.ServiceData), data.Temperature, data.Humidity, data.Battery)
	}
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("Done\n")
	case context.Canceled:
		fmt.Printf("Canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
