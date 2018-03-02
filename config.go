package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// RiemannConfig configuration for a Riemann output server
type RiemannConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	TLS      bool   `yaml:"tls"`
	KeyPath  string `yaml:"key_path"`
	CertPath string `yaml:"cert_path"`
	Insecure bool   `yaml:"insecure"`
}

// StrategyConfig configuration for an event forwarding strategy
type StrategyConfig struct {
	Type string `yaml:"type"`
}

// TCPConfig configuration for the Riemann Relay TCP Server
type TCPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	TLS      bool   `yaml:"tls"`
	KeyPath  string `yaml:"key_path"`
	CertPath string `yaml:"cert_path"`
}

// Config the global configuration for Riemann Relay
type Config struct {
	Riemann   []RiemannConfig `yaml:"riemann"`
	TCPServer TCPConfig       `yaml:"tcp"`
	Strategy  StrategyConfig  `yaml:"strategy"`
}

// GetConfig get Riemann relay configuration from a yaml file
func GetConfig(yamlPath string) (Config, error) {
	var config Config
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(yamlFile), &config)
	return config, err
}
