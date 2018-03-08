package strategy

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/mcorbin/riemann-relay/pkg/client"
	"github.com/mcorbin/riemann-relay/pkg/config"
	"github.com/riemann/riemann-go-client"
)

// Strategy an event forwarding Strategy.
type Strategy interface {
	Send(events *[]riemanngo.Event) []int
}

// BroadcastStrategy forward the events to ALL clients
type BroadcastStrategy struct {
	Clients []*client.Client
}

// Send send the events
func (s *BroadcastStrategy) Send(events *[]riemanngo.Event) []int {
	reconnectIndex := make([]int, 0)
	for i, client := range s.Clients {
		if client.Connected {
			result, err := riemanngo.SendEvents(client.Riemann, events)
			if err != nil {
				glog.Errorf("Error sending events: %s", err)
				err := client.Riemann.Close()
				client.Connected = false
				if err != nil {
					glog.Infof("Error closing connection: %s",
						err)
				}
				reconnectIndex = append(reconnectIndex, i)
			} else {
				glog.Info("Result: ", result)
			}
		} else {

			glog.Errorf("Error, this client is not connected: ", client.Config)
			reconnectIndex = append(reconnectIndex, i)
		}
	}
	return reconnectIndex
}

// GetStrategy takes the Strategy configuration, a slide of Client, and returns the
// event forwarding Strategy for Riemann Relay
func GetStrategy(config config.StrategyConfig, clients []*client.Client) (*BroadcastStrategy, error) {
	if config.Type == "broadcast" {
		strategy := &BroadcastStrategy{
			Clients: clients,
		}
		return strategy, nil
	}
	return nil, fmt.Errorf("Unknown strategy: %s", config.Type)
}
