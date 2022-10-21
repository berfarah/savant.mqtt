package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// Example Config:
// { "registry_file_path": "<path>.json",
//   "poll_seconds": 5
// }
type Config struct {
	RegistryFilePath string `json:"registry_file_path"`
	PollSeconds      int    `json:"poll_seconds"`
}

var configPath string
var registryPath string

const defaultConfigPath = "/usr/local/etc/savantmqtt.conf"
const defaultRegistryPath = "/usr/local/etc/savantregistry.json"

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "(required) Specifiy config file path")
	flag.StringVar(&registryPath, "r", "", "(optional) Specify savant registry path")
}

func LoadConfig() *Config {
	var config Config

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Couldn't open config file: %v\n", err)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal([]byte(bytes), &config); err != nil {
		log.Fatalf("Couldn't open config file: %v\n", err)
	}

	if registryPath != "" {
		config.RegistryFilePath = registryPath
	}

	if config.RegistryFilePath == "" {
		config.RegistryFilePath = defaultRegistryPath
	}

	if config.PollSeconds == 0 {
		config.PollSeconds = 5
	}

	return &config
}
