package main

import (
	"fmt"
	"net"
	"io"

	pb "github.com/golang/protobuf/proto"
	"github.com/riemann/riemann-go-client"
	"github.com/riemann/riemann-go-client/proto"
	"github.com/golang/glog"
)


// StartServer start a tcp server
func StartServer(addr string) error {
	glog.Info("Starting Riemann Relay TCP server...")
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		glog.Errorf("Error resolving TCP addr %s: %s", addr, err)
		return err
	}

	c := make(chan *[]riemanngo.Event)

	go func(){
		for{
			event := <-c
			fmt.Println(event)
		}
	}()

	listener, err := net.ListenTCP("tcp", address)

	if err != nil {
		glog.Errorf("Error creating TCP server on %s: %s ", addr, err)
		return err
	}

	glog.Info("Riemann Relay TCP server started on: ",
		listener.Addr().String())

	for {
		if conn, err := listener.AcceptTCP(); err == nil {
			go HandleConnection(conn, c)
		} else {
			glog.Error("Error accepting TCP connection ", err)
			return err
		}
	}
}

func getMsgSize(buffer []byte) byte {
	size := buffer[0] << 3 +
		buffer[1] << 2 +
		buffer[2] << 1 +
		buffer[3]
	return size
}

func getOkMsg() *proto.Msg {
	msg := new(proto.Msg)
	t := true
	msg.Ok = &t
	return msg
}


func getRespSizeBuffer(respBuffer *[]byte) []byte {
	size := len(*respBuffer)
	buffer := make([]byte, 4)
	buffer[0] = byte((size >> 3) & 0xE)
	buffer[1] = byte((size >> 2) & 0xE)
	buffer[2] = byte((size >> 1) & 0xE)
	buffer[3] = byte(size & 0xE)
	return buffer
}

func getErrorMsg(err error) *proto.Msg {
	msg := new(proto.Msg)
	f := false
	msg.Ok = &f
	errStr := err.Error()
	msg.Error = &errStr
	return msg
}

func writeError(conn net.Conn, err error) error {
	glog.Errorf("TCP error: %s", err.Error())
	msgBuffer, err := pb.Marshal(getErrorMsg(err))
	if err != nil {
		return err
	}
	msgSizeBuffer := getRespSizeBuffer(&msgBuffer)
	_,err = conn.Write(msgSizeBuffer)
	if err != nil {
		return err
	}
	_,err = conn.Write(msgBuffer)

	return err
}

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

// HandleConnection hande tcp connection
func HandleConnection(conn net.Conn, c chan *[]riemanngo.Event) {
	for {
		// read protobuf msg size
		sizeBuffer := make([]byte, 4)
		_,err := conn.Read(sizeBuffer)
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
		_,err = conn.Read(protoMsgBuffer)
		if err := checkTCPError(conn, err); err != nil {
			break
		}

		protoMsg := new(proto.Msg)
		err = pb.Unmarshal(protoMsgBuffer, protoMsg)
		if err := checkTCPError(conn, err); err != nil {
			break
		}

		events := riemanngo.ProtocolBuffersToEvents(protoMsg.Events)
		c <- &events

		msgBuffer, err := pb.Marshal(getOkMsg())

		if err := checkTCPError(conn, err); err != nil {
			break
		}

		msgSizeBuffer := getRespSizeBuffer(&msgBuffer)
		_,err = conn.Write(msgSizeBuffer)

		if err := checkTCPError(conn, err); err != nil {
			break
		}

		_,err = conn.Write(msgBuffer)
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