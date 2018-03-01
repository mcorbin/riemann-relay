package main

import (
	"github.com/riemann/riemann-go-client"
	"fmt"
)


func ConstructClients(config Config) ([]riemanngo.Client, error) {
	riemannConfig := config.Riemann
	clients := make([]riemanngo.Client, len(riemannConfig))
	for i, clientConfig := range riemannConfig {
		protocol := clientConfig.Protocol
		if protocol == "tcp" {
			tcpAddr := fmt.Sprintf("%s:%d", clientConfig.Host, clientConfig.Port)
			clients[i] = riemanngo.NewTcpClient(tcpAddr)
		} else {
			return nil, fmt.Errorf("Unknown protocol: %s", protocol)
		}
	}
	return clients, nil
}
