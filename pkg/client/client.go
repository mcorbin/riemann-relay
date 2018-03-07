package client

import (
	"fmt"
	"github.com/mcorbin/riemann-relay/pkg/config"
	"github.com/riemann/riemann-go-client"
)

// Client wrap the riemanngo.Client type
type Client struct {
	Riemann   riemanngo.Client
	Config    config.RiemannConfig
	Connected bool
}

// GetRiemannClient get a Riemann TCP client from a configuration
func GetRiemannClient(config config.RiemannConfig) riemanngo.Client {
	tcpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return riemanngo.NewTcpClient(tcpAddr)
}

// ConstructClients returns a slice of Client from a Riemann Relay configuration
func ConstructClients(config config.Config) ([]*Client, error) {
	riemannConfig := config.Riemann
	clients := make([]*Client, len(riemannConfig))
	for i, clientConfig := range riemannConfig {
		protocol := clientConfig.Protocol
		riemannClient := GetRiemannClient(clientConfig)
		if protocol == "tcp" {
			client := Client{
				Riemann: riemannClient,
				Config:  clientConfig,
			}
			clients[i] = &client
		}
	}
	return clients, nil
}
