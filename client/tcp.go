package client

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

type tcpTimeProtocolClient struct {
	addr    string
	timeout time.Duration
}

func newTCPTimeProtocolClient(addr string, timeout time.Duration) ITimeProtocolClient {
	return &tcpTimeProtocolClient{
		addr:    addr,
		timeout: timeout,
	}
}

func (c *tcpTimeProtocolClient) Get() (seconds uint32, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", c.addr); err != nil {
		return 0, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, err
	}

	var payload = make([]byte, 4)
	if _, err = conn.Read(payload); err != nil {
		return 0, err
	}

	var ts uint32
	buf := bytes.NewReader(payload)
	if err = binary.Read(buf, binary.BigEndian, &ts); err != nil {
		return 0, err
	}

	return ts, nil
}
