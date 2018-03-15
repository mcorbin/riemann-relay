package client

import (
	"github.com/mcorbin/riemann-relay/pkg/config"
	"github.com/riemann/riemann-go-client"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConstructClients(t *testing.T) {
	config := config.Config{
		Riemann: []config.RiemannConfig{
			config.RiemannConfig{
				Name:     "test",
				Host:     "localhost",
				Port:     5557,
				Protocol: "tcp",
			},
		},
	}
	clients, err := ConstructClients(config)
	assert.NoError(t, err)
	assert.Equal(t, len(clients), 1)
	clientType := reflect.TypeOf(clients[0].Riemann)
	assert.Equal(t, clientType, reflect.TypeOf(riemanngo.NewTcpClient("127.0.0.1:5555")))
}
