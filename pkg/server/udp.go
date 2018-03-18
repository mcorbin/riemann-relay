package server

// TODO server interface (start/stop... ) ?

import (
	"net"

	"github.com/golang/glog"
	pb "github.com/golang/protobuf/proto"
	"github.com/riemann/riemann-go-client"
	"github.com/riemann/riemann-go-client/proto"
)

// UDPServer a UDP server receiving Riemann events
type UDPServer struct {
	address   *net.UDPAddr
	stop      chan bool
	eventChan chan *[]riemanngo.Event
}

func NewUDPServer(addr string, c chan *[]riemanngo.Event) (*UDPServer, error) {
	address, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		glog.Errorf("Error resolving UDP addr %s: %s", addr, err)
		return nil, err
	}
	server := UDPServer{
		address:   address,
		stop:      make(chan bool),
		eventChan: c,
	}
	return &server, nil
}

func (server *UDPServer) StartServer() error {
	listener, err := net.ListenUDP("udp", server.address)
	if err != nil {
		glog.Errorf("Error creating UDP server on %s: %s ", server.address, err)
		return err
	}
	glog.Info("Riemann Relay UDP server started on: ",
		server.address.String())

	go func() {
		for {
			select {
			case <-server.stop:
				// close and stop server
				// err := conn.Close()
				// fmt.Println("close")
				// if err != nil {
				// 	glog.Errorf("UDP error during connection close ", err.Error())
				// }
			default:
				HandleUDPConnection(listener, server.eventChan)
			}

		}
	}()
	return nil
}

// HandleUDPConnection handle a new UDP connection to Riemann Relay
func HandleUDPConnection(conn *net.UDPConn, c chan *[]riemanngo.Event) {
	for {
		buffer := make([]byte, 16384)
		size, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			glog.Errorf("Error reading UDP datagram: %s",
				err.Error())
			break
		}
		protoMsg := new(proto.Msg)
		err = pb.Unmarshal(buffer[0:size], protoMsg)
		if err != nil {
			glog.Errorf("Error converting the UDP buffer into protobuf: %s",
				err.Error())
			break
		}
		events := riemanngo.ProtocolBuffersToEvents(protoMsg.Events)
		glog.Info(events)
		c <- &events
	}
}
