package client

import (
	"errors"
	"time"

	"timeProtocol/server"
)

type ITimeProtocolClient interface {
	Get() (seconds uint32, err error)
}

func NewTimeProtocolClient(proto server.Protocol, addr string, timeout time.Duration) (client ITimeProtocolClient, err error) {
	switch proto {
	case server.UDP:
		return newUDPTimeProtocolClient(addr, timeout), nil
	case server.TCP:
		return newTCPTimeProtocolClient(addr, timeout), nil
	}

	return nil, errors.New("protocol not supported")
}
