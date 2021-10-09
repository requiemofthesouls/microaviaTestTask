package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

type udpTimeProtocolClient struct {
	addr    string
	timeout time.Duration
}

func newUDPTimeProtocolClient(addr string, timeout time.Duration) ITimeProtocolClient {
	return &udpTimeProtocolClient{
		addr:    addr,
		timeout: timeout,
	}
}

func (c *udpTimeProtocolClient) Get() (seconds uint32, err error) {
	var conn net.Conn
	if conn, err = net.Dial("udp", c.addr); err != nil {
		return 0, err
	}

	if err = conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, err
	}

	if _, err = conn.Write(nil); err != nil {
		return 0, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, err
	}

	var payload = make([]byte, 4)
	if _, err = bufio.NewReader(conn).Read(payload); err != nil {
		return 0, err
	}

	var ts uint32
	buf := bytes.NewReader(payload)
	if err = binary.Read(buf, binary.BigEndian, &ts); err != nil {
		return 0, err
	}

	return ts, nil
}
