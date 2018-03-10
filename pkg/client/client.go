package client

import (
	"fmt"
	"github.com/mcorbin/riemann-relay/pkg/config"
	"github.com/mcorbin/riemann-relay/pkg/server"
	"github.com/riemann/riemann-go-client"
	"github.com/riemann/riemann-go-client/proto"
)

// Client wrap the riemanngo.Client type
type Client struct {
	Riemann   riemanngo.Client
	Config    config.RiemannConfig
	Connected bool
}

// GetRiemannClient get a Riemann TCP client from a configuration
func GetRiemannClient(config config.RiemannConfig) riemanngo.Client {
	tcpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return riemanngo.NewTcpClient(tcpAddr)
}

// ConstructClients returns a slice of Client from a Riemann Relay configuration
func ConstructClients(config config.Config) ([]*Client, error) {
	riemannConfig := config.Riemann
	clients := make([]*Client, len(riemannConfig))
	for i, clientConfig := range riemannConfig {
		protocol := clientConfig.Protocol
		riemannClient := GetRiemannClient(clientConfig)
		if protocol == "tcp" {
			client := Client{
				Riemann: riemannClient,
				Config:  clientConfig,
			}
			clients[i] = &client
		}
	}
	return clients, nil
}

// used in tests

type RiemannFixtureClient struct {
	messages chan *proto.Msg
}

func (c *RiemannFixtureClient) Connect(timeout int32) error {
	return nil
}

func (c *RiemannFixtureClient) Send(message *proto.Msg) (*proto.Msg, error) {
	go func() {
		c.messages <- message
	}()
	return server.NewOkMsg(), nil
}

func (t *RiemannFixtureClient) Close() error {
	return nil
}

func NewFixtureClient(sink chan *proto.Msg, connected bool) Client {
	client := Client{
		Riemann: &RiemannFixtureClient{
			messages: sink,
		},
		Config:    config.NewRiemannFixtureConfig(),
		Connected: connected,
	}
	return client
}

func NewFixtureClients(sink []chan *proto.Msg, connected bool) []*Client {
	clients := make([]*Client, len(sink))
	for i := range sink {
		c := NewFixtureClient(sink[i], connected)
		clients[i] = &c
	}

	return clients
}
