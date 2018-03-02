package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/golang/glog"
	pb "github.com/golang/protobuf/proto"
	"github.com/riemann/riemann-go-client"
	"github.com/riemann/riemann-go-client/proto"
)

// TCPServer a TCP server receiving Riemann events
type TCPServer struct {
	address *net.TCPAddr
	stop    chan bool
}

// StartServer start the Riemann Relay TCP server
func StartServer(addr string, c chan *[]riemanngo.Event) (*TCPServer, error) {
	glog.Info("Starting Riemann Relay TCP server...")
	address, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		glog.Errorf("Error resolving TCP addr %s: %s", addr, err)
		return nil, err
	}
	server := TCPServer{
		address: address,
		stop:    make(chan bool),
	}

	listener, err := net.ListenTCP("tcp", address)

	if err != nil {
		glog.Errorf("Error creating TCP server on %s: %s ", addr, err)
		return nil, err
	}

	glog.Info("Riemann Relay TCP server started on: ",
		listener.Addr().String())

	go func() {
		for {
			select {
			case <-server.stop:
				// stop server
			default:
				if conn, err := listener.AcceptTCP(); err == nil {
					go HandleConnection(conn, c)
				} else {
					glog.Error("Error accepting TCP connection ", err)
				}
			}
		}
	}()
	return &server, nil
}

// getMsgSize returns a uint32 from a slice of byte
func getMsgSize(buffer []byte) uint32 {
	return binary.BigEndian.Uint32(buffer)
}

// newOkMSg returns an OK proto.Msg
func newOkMsg() *proto.Msg {
	msg := new(proto.Msg)
	t := true
	msg.Ok = &t
	return msg
}

// newErrorMsg returns an Error proto.Msg
func newErrorMsg(err error) *proto.Msg {
	msg := new(proto.Msg)
	f := false
	msg.Ok = &f
	errStr := err.Error()
	msg.Error = &errStr
	return msg
}

// getRespSizeBuffer returns the size of a buffer, in a byte array (4 elements)
func getRespSizeBuffer(respBuffer []byte) []byte {
	size := len(respBuffer)
	buffer := make([]byte, 4)
	buffer[0] = byte((size >> 24) & 0xFF)
	buffer[1] = byte((size >> 16) & 0xFF)
	buffer[2] = byte((size >> 8) & 0xFF)
	buffer[3] = byte(size & 0xFF)
	return buffer
}

// writeError Takes a connection and an error, and send the error to the client
func writeError(conn net.Conn, err error) error {
	glog.Errorf("TCP error: %s", err.Error())
	msgBuffer, err := pb.Marshal(newErrorMsg(err))
	if err != nil {
		return err
	}
	msgSizeBuffer := getRespSizeBuffer(msgBuffer)
	_, err = conn.Write(msgSizeBuffer)
	if err != nil {
		return err
	}
	_, err = conn.Write(msgBuffer)

	return err
}

// checkTCPError take a connection and a optential error, send the error to
// the client if necessary
func checkTCPError(conn net.Conn, err error) error {
	if err != nil {
		if err := writeError(conn, err); err != nil {
			glog.Errorf("TCP unrecoverable error: %s",
				err.Error())
			return err
		}
	}
	return nil
}

// HandleConnection handle a new TCP connection to Riemann Relay
func HandleConnection(conn net.Conn, c chan *[]riemanngo.Event) {
	for {
		// read protobuf msg size
		sizeBuffer := make([]byte, 4)
		_, err := conn.Read(sizeBuffer)
		if err != nil {
			if err != io.EOF {
				if err := checkTCPError(conn, err); err != nil {
					break
				}
			} else {
				// close connection if EOF
				break
			}
		}

		msgSize := getMsgSize(sizeBuffer)
		// read protobuf msg
		protoMsgBuffer := make([]byte, msgSize)
		_, err = conn.Read(protoMsgBuffer)
		if err := checkTCPError(conn, err); err != nil {
			break
		}

		protoMsg := new(proto.Msg)
		err = pb.Unmarshal(protoMsgBuffer, protoMsg)
		if err := checkTCPError(conn, err); err != nil {
			break
		}
		events := riemanngo.ProtocolBuffersToEvents(protoMsg.Events)
		glog.Info(events)
		c <- &events
		msgBuffer, err := pb.Marshal(newOkMsg())

		if err := checkTCPError(conn, err); err != nil {
			break
		}

		msgSizeBuffer := getRespSizeBuffer(msgBuffer)
		_, err = conn.Write(msgSizeBuffer)

		if err := checkTCPError(conn, err); err != nil {
			break
		}

		_, err = conn.Write(msgBuffer)
		if err := checkTCPError(conn, err); err != nil {
			break
		}
	}

	err := conn.Close()
	fmt.Println("close")
	if err != nil {
		glog.Errorf("TCP error during connection close ", err.Error())
	}
}
