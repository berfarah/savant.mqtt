package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// Example Config:
// { "topic_prefix": "homeassistant",
//   "topic_node_id": "savant",
//   "broker": "<tcp|ssl>://<url>:<port>",
//   "ssl_cert_path": "<path>.crt",
//   "registry_file_path": "<path>.json",
//   "poll_seconds": 5
// }
type Config struct {
	TopicPrefix      string `json:"topic_prefix"`
	TopicNodeID      string `json:"topic_node_id"`
	Broker           string `json:"broker"`
	UseSSL           bool   `json:"-"`
	SSLCertPath      string `json:"ssl_cert_path"`
	RegistryFilePath string `json:"registry_file_path"`
	PollSeconds      int    `json:"poll_seconds"`
	Username         string `json:"username"`
	Password         string `json:"password"`
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

	if config.TopicPrefix == "" {
		config.TopicPrefix = "homeassistant"
	}

	if config.TopicNodeID == "" {
		config.TopicNodeID = "savant"
	}

	if config.SSLCertPath != "" {
		config.UseSSL = true
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
