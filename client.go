package main

import (
	"fmt"
	"github.com/riemann/riemann-go-client"
)

type ClientWrapper struct {
	client riemanngo.Client
	config RiemannConfig
}

func GetRiemannClient(config RiemannConfig) riemanngo.Client {
	tcpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return riemanngo.NewTcpClient(tcpAddr)
}

func ConstructClient(config RiemannConfig) (*ClientWrapper, error) {
	protocol := config.Protocol
	if protocol == "tcp" {
		tcpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		client := ClientWrapper{
			client: riemanngo.NewTcpClient(tcpAddr),
			config: config,
		}
		return &client, nil
	} else {
		return nil, fmt.Errorf("Unknown protocol: %s", protocol)
	}
}

func ConstructClients(config Config) ([]*ClientWrapper, error) {
	riemannConfig := config.Riemann
	clients := make([]*ClientWrapper, len(riemannConfig))
	for i, clientConfig := range riemannConfig {
		protocol := clientConfig.Protocol
		riemannClient := GetRiemannClient(clientConfig)
		if protocol == "tcp" {
			client := ClientWrapper{
				client: riemannClient,
				config: clientConfig,
			}
			clients[i] = &client
		}
	}
	return clients, nil
}
