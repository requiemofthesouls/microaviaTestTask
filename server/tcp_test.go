package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"timeProtocol/service"
)

const TCPAddr = ":57346"

var timeout = time.Duration(10) * time.Second
var TCPServer ITimeProtocolServer
var svc = service.NewTimeProtocolService(clock.New())

func tearUpTCP() {
	var err error

	if TCPServer, err = NewTimeServer(TCP, TCPAddr, 1, timeout, svc); err != nil {
		log.Fatal(err)
	}

	if err = TCPServer.Run(); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	go tearUpTCP()
	go tearUpUDP()
	time.Sleep(time.Second)
	code := m.Run()
	TCPServer.Stop()
	UDPServer.Stop()
	time.Sleep(time.Second)
	os.Exit(code)
}

func TestGetTimeTCP(t *testing.T) {
	var err error
	var conn net.Conn
	if conn, err = net.Dial("tcp", TCPAddr); err != nil {
		t.Error(err)
	}

	var payload = make([]byte, 4)
	if _, err = conn.Read(payload); err != nil {
		t.Error(err)
	}

	var ts uint32
	buf := bytes.NewReader(payload)
	if err = binary.Read(buf, binary.BigEndian, &ts); err != nil {
		t.Error(err)
	}

	log.Println("raw response:", ts)

	dur := time.Duration(ts) * time.Second
	gotTimeNow := service.TPStartTime.Add(dur)
	expectedTimeNow := time.Now().In(time.UTC).Truncate(time.Second)

	log.Println("got time:", gotTimeNow)
	log.Println("expected time:", expectedTimeNow)

	assert.Equal(t, gotTimeNow, expectedTimeNow)
}
