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
		glog.Error("Error: ", err)
		os.Exit(1)
	}
	
	c := make(chan *[]riemanngo.Event)
	go func() {
		for {
			event := <-c
			fmt.Println(event)
		}
	}()


	tcpAddr := fmt.Sprintf("%s:%d", config.TCPServer.Host, config.TCPServer.Port)
	_, err = StartServer(tcpAddr, c)
	if err != nil {
		glog.Errorf("Stopping Riemann Relay: %s", err.Error())
	}
	select {}
}
