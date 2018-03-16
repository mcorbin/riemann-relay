package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/mcorbin/riemann-relay/pkg/client"
	"github.com/mcorbin/riemann-relay/pkg/config"
	"github.com/mcorbin/riemann-relay/pkg/server"
	"github.com/mcorbin/riemann-relay/pkg/strategy"
	"github.com/riemann/riemann-go-client"
)

func main() {
	configPath := flag.String("config", "", "Riemann Relay config path")
	flag.Parse()
	if *configPath == "" {
		glog.Error("--config missing")
		os.Exit(1)
	}
	config, err := config.GetConfig(*configPath)
	if err != nil {
		glog.Error("Error loading the configuration: ", err)
		os.Exit(1)
	}
	clients, err := client.ConstructClients(config)
	if err != nil {
		glog.Error("Error creating the Riemann clients: ", err)
		os.Exit(1)
	}
	strategy, err := strategy.GetStrategy(config.Strategy, clients)
	if err != nil {
		glog.Error("Error creating strategy: ", err)
		os.Exit(1)
	}

	c := make(chan *[]riemanngo.Event)
	for _, client := range clients {
		err := client.Riemann.Connect(5)
		if err != nil {
			glog.Errorf("Error connecting %s: %s", client.Config, err)
		} else {
			client.Connected = true
		}
	}

	go func() {
		for {
			events := <-c
			reconnectIndex := strategy.Send(events)
			// reconnect
			for _, i := range reconnectIndex {
				glog.Info("Trying to reconnect ")
				config := strategy.Clients[i].Config
				client := client.GetRiemannClient(config)
				err := client.Connect(5)
				if err != nil {
					glog.Errorf("Reconnect connect %s failed: %s",
						config,
						err)
				} else {
					strategy.Clients[i].Riemann = client
					strategy.Clients[i].Connected = true
					glog.Infof("Connected again ! %s", config)
				}

			}
		}
	}()

	tcpAddr := fmt.Sprintf("%s:%d", config.TCPServer.Host, config.TCPServer.Port)
	tcpServer, err := server.NewTCPServer(tcpAddr, c)
	if err != nil {
		// TODO better error handling/msg
		glog.Errorf("Stopping Riemann Relay: %s", err.Error())
	}
	err = tcpServer.StartServer()
	if err != nil {
		glog.Errorf("Stopping Riemann Relay: %s", err.Error())
	}
	select {}
}
