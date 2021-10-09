package service

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
)

var (
	TPStartTime     = time.Date(1900, 01, 01, 00, 00, 00, 00, time.UTC)
	SystemTimeError = fmt.Errorf("system time is less than %s", TPStartTime)
)

type ITimeProtocolService interface {
	GetBinarySeconds() []byte
	GetSeconds() uint32
}

type service struct {
	clock clock.Clock
}

func NewTimeProtocolService(clock clock.Clock) ITimeProtocolService {
	return &service{clock: clock}
}

func (s service) GetBinarySeconds() []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, s.GetSeconds())
	return buf
}

func (s service) GetSeconds() uint32 {
	now := s.clock.Now()
	if now.Before(TPStartTime) {
		panic(SystemTimeError)
	}
	return uint32(s.clock.Since(TPStartTime).Seconds())
}
