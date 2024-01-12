package btle

import (
	"bytes"

	"github.com/go-ble/ble"
)

type DevData struct {
	ServiceData      []byte
	ManufacturerData []byte
	Temperature      float64
	Humidity         float64
	Battery          float64
}

var (
	BTDevice = make(map[string]DevData)
)

func Handler(a ble.Advertisement) {
	if len(a.ManufacturerData()) > 0 {
		if bytes.Equal(a.ManufacturerData()[0:2], []byte{0x69, 0x09}) {
			_, ok := BTDevice[a.Addr().String()]
			if !ok {
				BTDevice[a.Addr().String()] = DevData{ManufacturerData: []byte{}, ServiceData: []byte{}}
			}
			if len(a.ManufacturerData()) > 0 {
				if entry, ok := BTDevice[a.Addr().String()]; ok {
					entry.ManufacturerData = a.ManufacturerData()
					entry.Temperature = buildTemperature(entry.ManufacturerData)
					entry.Humidity = buildHumidity(entry.ManufacturerData)
					BTDevice[a.Addr().String()] = entry
				}
			}
			if len(a.ServiceData()) > 0 {
				if entry, ok := BTDevice[a.Addr().String()]; ok {
					entry.ServiceData = a.ServiceData()[0].Data
					entry.Battery = buildBattery(a.ServiceData()[0].Data)
					BTDevice[a.Addr().String()] = entry
				}
			}
		}
	}
}

func buildTemperature(t []byte) float64 {
	temperature := (float64(t[10]&0x0f)*0.1 + float64(t[11]&0x7f))
	sign := 1
	if (t[11] & 0x80) == 0 {
		sign = -1
	}
	return temperature * float64(sign)
}

func buildHumidity(h []byte) float64 {
	humidity := float64(h[12] & 0x7f)
	return humidity
}

func buildBattery(b []byte) float64 {
	return float64(b[2] & 0x7f)
}
