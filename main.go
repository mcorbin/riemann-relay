package main

import (
	"fmt"
	"github.com/riemann/riemann-go-client"
	"github.com/golang/glog"
)

func main() {
	c := make(chan *[]riemanngo.Event)
	go func(){
		for{
			event := <-c
			fmt.Println(event)
		}
	}()

	err := StartServer("127.0.0.1:2124", c)
	if err != nil {
		glog.Errorf("Stopping Riemann Relay: %s", err.Error())
	}
}
