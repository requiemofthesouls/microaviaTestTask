package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"timeProtocol/service"
)

const UDPAddr = ":57347"

var UDPServer ITimeProtocolServer

func tearUpUDP() {
	var err error

	if UDPServer, err = NewTimeServer(UDP, UDPAddr, 1, timeout, svc); err != nil {
		log.Fatal(err)
	}

	if err = UDPServer.Run(); err != nil {
		log.Fatal(err)
	}
}

func TestNewUDPTimeServer(t *testing.T) {
	udpAddr, err := net.ResolveUDPAddr("udp4", UDPAddr)
	assert.NoError(t, err)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	assert.NoError(t, err)

	_, err = conn.Write([]byte("anything"))
	assert.NoError(t, err)

	var buf [512]byte
	n, err := conn.Read(buf[0:])
	assert.NoError(t, err)
	assert.Equal(t, 4, n)

	var ts uint32
	err = binary.Read(bytes.NewReader(buf[:n]), binary.BigEndian, &ts)
	assert.NoError(t, err)

	log.Println("raw response:", ts)

	dur := time.Duration(ts) * time.Second
	gotTimeNow := service.TPStartTime.Add(dur)
	expectedTimeNow := time.Now().In(time.UTC).Truncate(time.Second)

	log.Println("got time:", gotTimeNow)
	log.Println("expected time:", expectedTimeNow)

	assert.Equal(t, gotTimeNow, expectedTimeNow)
}
