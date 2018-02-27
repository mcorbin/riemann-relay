package main


import (
	"testing"
	"github.com/stretchr/testify/assert"
)


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
