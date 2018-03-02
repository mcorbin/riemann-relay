package main

import (
	"fmt"
	"github.com/riemann/riemann-go-client"
)

// Client wrap the riemanngo.Client type
type Client struct {
	riemann   riemanngo.Client
	config    RiemannConfig
	connected bool
}

// GetRiemannClient get a Riemann TCP client from a configuration
func GetRiemannClient(config RiemannConfig) riemanngo.Client {
	tcpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return riemanngo.NewTcpClient(tcpAddr)
}

// ConstructClients returns a slice of Client from a Riemann Relay configuration
func ConstructClients(config Config) ([]*Client, error) {
	riemannConfig := config.Riemann
	clients := make([]*Client, len(riemannConfig))
	for i, clientConfig := range riemannConfig {
		protocol := clientConfig.Protocol
		riemannClient := GetRiemannClient(clientConfig)
		if protocol == "tcp" {
			client := Client{
				riemann: riemannClient,
				config:  clientConfig,
			}
			clients[i] = &client
		}
	}
	return clients, nil
}
