package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	config, err := GetConfig("config_test.yml")
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", config.TCPServer.Host)
	assert.Equal(t, 2120, config.TCPServer.Port)
	assert.Equal(t, 2, len(config.Riemann))
	assert.Equal(t, "local", config.Riemann[0].Name)
	assert.Equal(t, "localhost", config.Riemann[0].Host)
	assert.Equal(t, 5555, config.Riemann[0].Port)
	assert.Equal(t, "tcp", config.Riemann[0].Protocol)
}
