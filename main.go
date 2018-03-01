package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/riemann/riemann-go-client"
)

func main() {
	configPath := flag.String("config", "", "Riemann Relay config path")
	flag.Parse()
	if *configPath == "" {
		glog.Error("--config missing")
		os.Exit(1)
	}
	config, err := GetConfig(*configPath)
	if err != nil {
		glog.Error("Error loading the configuration: ", err)
		os.Exit(1)
	}
	clients, err := ConstructClients(config)
	if err != nil {
		glog.Error("Error creating the Riemann clients: ", err)
		os.Exit(1)
	}
	c := make(chan *[]riemanngo.Event)
	for _, client := range clients {
		client.Connect(5)
	}

	go func() {
		for {
			events := <-c
			for _, client := range clients {
				result, err := riemanngo.SendEvents(client, events)
				if err != nil {
					glog.Errorf("Error sending events: %s", err)
				} else {
					glog.Info("Result: ", result)
				}
			}
		}
	}()


	tcpAddr := fmt.Sprintf("%s:%d", config.TCPServer.Host, config.TCPServer.Port)
	_, err = StartServer(tcpAddr, c)
	if err != nil {
		glog.Errorf("Stopping Riemann Relay: %s", err.Error())
	}
	select {}
}
