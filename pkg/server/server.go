package server

import (
	"github.com/riemann/riemann-go-client"
)

// Server the server interface
type Server interface {
	Send(events *[]riemanngo.Event) []int
	StartServer(addr string, c chan *[]riemanngo.Event) (*Server, error)
}
