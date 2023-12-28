package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Cfg struct {
	Sensors map[string]Sensor
}

type Sensor struct {
	Labels []Label
}

type Label struct {
	Name  string
	Value string
}

var (
	Config         *Cfg
	ConfigFile     string
	Device         string
	BTScanDuration time.Duration
	BTScanInterval time.Duration
)

func Init(configFile string) error {
	ConfigFile = configFile
	content := []byte(`{}`)
	_, err := os.Stat(configFile)
	if !os.IsNotExist(err) {
		content, err = os.ReadFile(configFile)
		if err != nil {
			return err
		}
	}
	if len(content) == 0 {
		content = []byte(`{}`)
	}
	err = json.Unmarshal(content, &Config)
	if err != nil {
		return err
	}

	log.Printf("Finished reading config.")
	return nil
}
