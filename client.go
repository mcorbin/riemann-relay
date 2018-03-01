package main

import (
	"github.com/riemann/riemann-go-client"
	"fmt"
)

type Client struct {
	client riemanngo.Client
	config RiemannConfig
	reconnect bool
}


func GetRiemannClient(config RiemannConfig) riemanngo.Client {
	tcpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return riemanngo.NewTcpClient(tcpAddr)
}

func ConstructClient(config RiemannConfig) (*Client, error) {
	protocol := config.Protocol
	if protocol == "tcp" {
		tcpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		client := Client {
			client: riemanngo.NewTcpClient(tcpAddr),
			config: config,
			reconnect: false,
		}
		return &client, nil
	} else {
		return nil, fmt.Errorf("Unknown protocol: %s", protocol)
	}
}


func ConstructClients(config Config) ([]*Client, error) {
	riemannConfig := config.Riemann
	clients := make([]*Client, len(riemannConfig))
	for i, clientConfig := range riemannConfig {
		protocol := clientConfig.Protocol
		riemannClient := GetRiemannClient(clientConfig)
		if protocol == "tcp" {
			client := Client {
				client: riemannClient,
				config: clientConfig,
				reconnect: false,
			}
			clients[i] = &client
		}
	}
	return clients, nil
}
