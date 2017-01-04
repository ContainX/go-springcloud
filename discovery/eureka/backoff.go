package eureka

import (
	"github.com/cenkalti/backoff"
	"time"
)

type MaxAttemptBackoff struct {
	Interval time.Duration
	Attempts int
	count    int
}

func NewMaxAttemptBackoff(interval time.Duration, attempts int) *MaxAttemptBackoff {
	return &MaxAttemptBackoff{Interval: interval, Attempts: attempts}
}

func (b *MaxAttemptBackoff) Reset() {
	b.count = 0
}

func (b *MaxAttemptBackoff) NextBackOff() time.Duration {
	if b.count >= (b.Attempts - 1) {
		b.Reset()
		return backoff.Stop
	}
	b.count = b.count + 1
	return b.Interval
}
