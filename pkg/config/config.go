package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

// TCPConfig configuration for the Riemann Relay TCP Server
type UDPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Config the global configuration for Riemann Relay
type Config struct {
	Riemann  []RiemannConfig `yaml:"riemann"`
	TCP      TCPConfig       `yaml:"tcp"`
	UDP      UDPConfig       `yaml:"udp"`
	Strategy StrategyConfig  `yaml:"strategy"`
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

// used in tests

func NewRiemannFixtureConfig() RiemannConfig {
	return RiemannConfig{
		Name:     "test",
		Host:     "localhost",
		Port:     5557,
		Protocol: "test",
	}
}

func NewFixtureConfig(s string) Config {
	config := Config{
		Riemann: []RiemannConfig{
			NewRiemannFixtureConfig(),
		},
		TCP: TCPConfig{
			Host: "localhost",
			Port: 2120,
		},
		Strategy: StrategyConfig{
			Type: s,
		},
	}
	return config
}
