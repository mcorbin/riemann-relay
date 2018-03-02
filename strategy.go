package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/riemann/riemann-go-client"
)

type Strategy interface {
	Send(events *[]riemanngo.Event)
	Reconnect(clientsWrapper []*ClientWrapper, reconnectIndex []int)
}

type BroadcastStrategy struct {
	clientsWrapper []*ClientWrapper
}

func (s *BroadcastStrategy) Send(events *[]riemanngo.Event) {
	reconnectIndex := make([]int, 0)
	for i, client := range s.clientsWrapper {
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

	// reconnect
	for _, i := range reconnectIndex {
		glog.Info("Trying to reconnect ")
		config := s.clientsWrapper[i].config
		client := GetRiemannClient(config)
		err := client.Connect(5)
		if err != nil {
			glog.Errorf("Reconnect connect %s failed: %s",
				config,
				err)
		} else {
			s.clientsWrapper[i].client = client
			glog.Infof("Connected again ! %s", config)
		}

	}
}

func GetStrategy(config StrategyConfig, clientsWrapper []*ClientWrapper) (*BroadcastStrategy, error) {
	if config.Type == "broadcast" {
		strategy := &BroadcastStrategy{
			clientsWrapper: clientsWrapper,
		}
		return strategy, nil
	}
	return nil, fmt.Errorf("Unknown strategy: %s", config.Type)
}
