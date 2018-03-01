package main

import (
	"testing"
	"time"

	"github.com/riemann/riemann-go-client"
	"github.com/stretchr/testify/assert"
)

// func init() {
// 	flag.Set("alsologtostderr", fmt.Sprintf("%t", true))
// 	var logLevel string
// 	flag.StringVar(&logLevel, "logLevel", "4", "test")
// 	flag.Lookup("v").Value.Set(logLevel)
// }

func TestGetMsgSize(t *testing.T) {
	assert.Equal(t, getMsgSize([]byte{0, 0, 0, 0}), uint32(0))
	assert.Equal(t, getMsgSize([]byte{0, 0, 0, 1}), uint32(1))
	assert.Equal(t, getMsgSize([]byte{0, 0, 0, 255}), uint32(255))
	assert.Equal(t, getMsgSize([]byte{255, 255, 255, 255}), uint32(0xFFFFFFFF))
}

func TestGetRespSizeBuffer(t *testing.T) {
	b := [5]byte{}
	assert.Equal(t, getRespSizeBuffer(b[0:4]), []byte{0, 0, 0, 4})
	c := [100000]byte{}
	assert.Equal(t, getRespSizeBuffer(c[0:100000]), []byte{0, 1, 0x86, 0xA0})
}

func TestTcpServer(t *testing.T) {
	c := make(chan *[]riemanngo.Event)
	_, err := StartServer("127.0.0.1:2120", c)
	assert.NoError(t, err)
	client := riemanngo.NewTcpClient("127.0.0.1:2120")

	client.Connect(5)

	now := time.Now().Round(time.Microsecond)

	event := riemanngo.Event{
		Service: "hello",
		Metric:  int64(100),
		Time:    now,
	}
	t.Log("time = ", now)
	go func() {

		assert.Equal(t, <-c, []riemanngo.Event{event})
	}()

	result, err := riemanngo.SendEvent(client, &event)
	assert.NoError(t, err)
	assert.Equal(t, result, newOkMsg())
}
