package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"timeProtocol/service"
)

func TestNewTimeServer(t *testing.T) {
	type args struct {
		proto   Protocol
		addr    string
		workers int
		timeout time.Duration
		svc     service.ITimeProtocolService
	}
	argsTCP := args{
		proto:   TCP,
		addr:    ":8080",
		workers: 1,
		timeout: timeout,
		svc:     svc,
	}

	argsUDP := args{
		proto:   UDP,
		addr:    ":8080",
		workers: 1,
		timeout: timeout,
		svc:     svc,
	}

	argsUnknown := args{
		proto:   124,
		addr:    ":8080",
		workers: 1,
		timeout: 1,
		svc:     svc,
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		errExpected error
	}{
		{
			name: "build TCP TCPServer",
			args: argsTCP,
		},
		{
			name: "build UDP TCPServer",
			args: argsUDP,
		},
		{
			name:        "unknown protocol error",
			args:        argsUnknown,
			wantErr:     true,
			errExpected: ProtocolNotSupportedError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTimeServer(tt.args.proto, tt.args.addr, tt.args.workers, tt.args.timeout, tt.args.svc)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errExpected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
