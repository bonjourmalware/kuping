package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/c2h5oh/datasize"
	"github.com/go-yaml/yaml"
)

var (
	Cfg = new(Config)
)

type ServerConfig struct {
	Enabled    bool     `yaml:"enabled"`
	PortsRules []string `yaml:"ports"`
	HostsRules []string `yaml:"hosts"`
	IPRules    IPRules

	EnabledPorts EnabledPorts
}

// Config structure which mirrors the yaml file
type Config struct {
	Logfile            string       `yaml:"Logfile"`
	HTTP               ServerConfig `yaml:"HTTP"`
	HTTPS              ServerConfig `yaml:"HTTPS"`
	MaxPOSTDataSizeRaw string       `yaml:"MaxPOSTDataSize"`
	EnableBlacklist    bool         `yaml:"EnableBlacklist"`
	EnableWhitelist    bool         `yaml:"EnableWhitelist"`
	ServerHeader       string       `yaml:"ServerHeader"`

	// Not parsed from config file
	MaxPOSTDataSize uint64
}

func (cfg *Config) Load() {
	var byteSize datasize.ByteSize

	filepath := "config.yml"
	cfgData, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Failed to load config file", err)
	}

	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		fmt.Printf("Failed to load the config file (%s)\n", filepath)
		fmt.Println(err)
		os.Exit(1)
	}

	if err := byteSize.UnmarshalText([]byte(cfg.MaxPOSTDataSizeRaw)); err != nil {
		fmt.Printf("Failed to parse the MaxPOSTDataSize value (%s)\n", cfg.MaxPOSTDataSizeRaw)
		fmt.Println(err)
		os.Exit(1)
	}

	if cfg.ServerHeader == "" {
		cfg.ServerHeader = "Apache"
	}

	cfg.MaxPOSTDataSize = byteSize.Bytes()

	cfg.HTTP.EnabledPorts.ParseRules(cfg.HTTP.PortsRules)
	cfg.HTTPS.EnabledPorts.ParseRules(cfg.HTTPS.PortsRules)

	cfg.HTTP.IPRules.ParseRules(cfg.HTTP.HostsRules)
	cfg.HTTPS.IPRules.ParseRules(cfg.HTTPS.HostsRules)
}
