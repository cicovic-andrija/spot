package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// DBConfig is a database configuration object
type DBConfig struct {
	ConnString string `json:"conn_string"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

// Config is a configuration object
type Config struct {
	Version   string   `json:"version"`
	AssetsDir string   `json:"assets_dir"`
	DevAddr   string   `json:"dev_addr"`
	DevPort   int      `json:"dev_port"`
	DBConfig  DBConfig `json:"db_config"`
}

// ReadConfig reads a configuration file
func ReadConfig(fileName string) (cfg Config, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		err = fmt.Errorf("Failed to open config file %s: %v", fileName, err)
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		err = fmt.Errorf("File %s decoding error: %v", fileName, err)
		return
	}

	return
}
