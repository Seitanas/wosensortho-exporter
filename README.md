# wosensortho-exporter

Exporter designed to collect and build metrics for `SwitchBot Indoor/Outdoor Thermo-Hygrometer` also known as `WoSensorTHO` or `Model W3400010`.

To start collecting metrics you need to know MAC address of your sensor device.
You can use scanner application to scan for available sensor devices:

```
go run cmd/scanner/main.go
Scanning for 5s...
Done
Found SwitchBot devices:
MAC: 11:11:11:00:bb:aa, ManufacturerData: 690911111100bbaa920308972800 ServiceData: 7700e4
```
Having found MAC addreses you can configure exporter service. Example config is presented in `config.json` file.
