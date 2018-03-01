package main

import (
	"github.com/riemann/riemann-go-client"
	"github.com/golang/glog"
	"fmt"
)

type Strategy interface {
	Send(events *[]riemanngo.Event)
}

type BroadcastStrategy struct {
	clients []*Client
}

func (s *BroadcastStrategy) Send(events *[]riemanngo.Event) {
	reconnectIndex := make([]int, 0)
	for i, client := range s.clients {
		glog.Info("Reconnect :", client.reconnect)
		result, err := riemanngo.SendEvents(client.client, events)
		if err != nil {
			glog.Errorf("Error sending events: %s", err)
			err := client.client.Close()
			if err != nil {
			glog.Infof("Error closing connection: %s",
				err)
			}
			reconnectIndex = append(reconnectIndex, i)
		} else {
			glog.Info("Result: ", result)
		}
	}

	for _,i := range reconnectIndex {
		glog.Info("Trying to reconnect ")
		config := s.clients[i].config
		client := GetRiemannClient(config)
		err := client.Connect(5)
		if err != nil {
			glog.Errorf("Reconnect connect %s failed: %s",
				config,
				err)
		} else {
			s.clients[i].client = client
			glog.Infof("Connected again ! %s", config)
		}

	}
}

func GetStrategy(config StrategyConfig, clients []*Client) (*BroadcastStrategy, error){
	if config.Type == "broadcast" {
		strategy := &BroadcastStrategy{
			clients: clients,
		}
		return strategy, nil
	}
	return nil, fmt.Errorf("Unknown strategy: %s", config.Type)
}
