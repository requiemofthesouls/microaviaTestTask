package service

import (
	"math/rand"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func Test_service_GetBinarySeconds(t *testing.T) {
	m := clock.NewMock()
	svc := NewTimeProtocolService(m)

	t.Run("positive scenario", func(t *testing.T) {
		t.Run("fist 1000 seconds", func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				m.Set(TPStartTime.Add(time.Second * time.Duration(i)))
				got := int(svc.GetSeconds())
				assert.Equal(t, got, i)
			}
		})

		t.Run("random 1000 seconds", func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				expected := rand.Uint32()
				m.Set(TPStartTime.Add(time.Second * time.Duration(expected)))
				got := svc.GetSeconds()
				assert.Equal(t, got, expected)
			}
		})
	})

	t.Run("negative scenario", func(t *testing.T) {
		t.Run("system date is less than 1900 year", func(t *testing.T) {
			startTime := time.Date(1800, 01, 01, 00, 00, 00, 00, time.UTC)
			m.Set(startTime)
			assert.PanicsWithError(t, SystemTimeError.Error(), func() {
				_ = svc.GetSeconds()
			})
		})
	})
}
