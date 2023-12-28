package btle

import (
	"bytes"

	"github.com/go-ble/ble"
)

type DevData struct {
	ServiceData      []byte
	ManufacturerData []byte
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
					BTDevice[a.Addr().String()] = entry
				}
			}
			if len(a.ServiceData()) > 0 {
				if entry, ok := BTDevice[a.Addr().String()]; ok {
					entry.ServiceData = a.ServiceData()[0].Data
					BTDevice[a.Addr().String()] = entry
				}
			}
		}
	}
}
