package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/riemann/riemann-go-client"
)

// Strategy an event forwarding Strategy.
type Strategy interface {
	Send(events *[]riemanngo.Event)
}

// BroadcastStrategy forward the events to ALL clients
type BroadcastStrategy struct {
	clients []*Client
}

// Send send the events
func (s *BroadcastStrategy) Send(events *[]riemanngo.Event) {
	reconnectIndex := make([]int, 0)
	for i, client := range s.clients {
		if client.connected {
			result, err := riemanngo.SendEvents(client.riemann, events)
			if err != nil {
				glog.Errorf("Error sending events: %s", err)
				err := client.riemann.Close()
				client.connected = false
				if err != nil {
					glog.Infof("Error closing connection: %s",
						err)
				}
				reconnectIndex = append(reconnectIndex, i)
			} else {
				glog.Info("Result: ", result)
			}
		} else {

			glog.Errorf("Error, this client is not connected: ", client.config)
			reconnectIndex = append(reconnectIndex, i)
		}
	}

	// reconnect
	for _, i := range reconnectIndex {
		glog.Info("Trying to reconnect ")
		config := s.clients[i].config
		client := GetRiemannClient(config)
		err := client.Connect(5)
		if err != nil {
			glog.Errorf("Reconnect connect %s failed: %s",
				config,
				err)
		} else {
			s.clients[i].connected = true
			s.clients[i].riemann = client
			glog.Infof("Connected again ! %s", config)
		}

	}
}

// GetStrategy takes the Strategy configuration, a slide of Client, and returns the
// event forwarding Strategy for Riemann Relay
func GetStrategy(config StrategyConfig, clients []*Client) (*BroadcastStrategy, error) {
	if config.Type == "broadcast" {
		strategy := &BroadcastStrategy{
			clients: clients,
		}
		return strategy, nil
	}
	return nil, fmt.Errorf("Unknown strategy: %s", config.Type)
}
