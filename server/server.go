package server

import (
	"errors"
	"time"

	"timeProtocol/service"
)

type (
	ITimeProtocolServer interface {
		Run() (err error)
		Stop()
	}

	Protocol int
)

const (
	UDP Protocol = iota
	TCP
)

var ProtocolNotSupportedError = errors.New("protocol not supported")

func NewTimeServer(
	proto Protocol,
	addr string,
	workers int,
	timeout time.Duration,
	svc service.ITimeProtocolService,
) (server ITimeProtocolServer, err error) {
	switch proto {
	case UDP:
		return newUDPTimeServer(addr, workers, timeout, svc), nil
	case TCP:
		return newTCPTimeServer(addr, workers, timeout, svc), nil
	}

	return nil, ProtocolNotSupportedError
}
