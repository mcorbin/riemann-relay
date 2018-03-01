package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)


type RiemannConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Protocol string `yaml:"protocol"`
	TLS bool `yaml:"tls"`
	KeyPath string `yaml:"key_path"`
	CertPath string `yaml:"cert_path"`
	Insecure bool `yaml:"insecure"`
}

type StrategyConfig struct {
	Type string `yaml:"type"`
}

type TCPConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	TLS bool `yaml:"tls"`
	KeyPath string `yaml:"key_path"`
	CertPath string `yaml:"cert_path"`
}

type Config struct {
	Riemann []RiemannConfig `yaml:"riemann"`
	TCPServer TCPConfig `yaml:"tcp"`
	Strategy StrategyConfig `yaml:"strategy"`
}

// GetConfig get Riemann relay configuration from yaml
func GetConfig(yamlPath string) (Config, error) {
	var config Config
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
