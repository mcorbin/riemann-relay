package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/riemann/riemann-go-client"
)

func main() {
	c := make(chan *[]riemanngo.Event)
	go func() {
		for {
			event := <-c
			fmt.Println(event)
		}
	}()

	_, err := StartServer("127.0.0.1:2124", c)
	if err != nil {
		glog.Errorf("Stopping Riemann Relay: %s", err.Error())
	}
	select {}
}
