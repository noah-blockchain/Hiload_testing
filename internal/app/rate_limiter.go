package app

import (
	"fmt"
	"math"
	"time"
)

type RateLimiter struct {
	Freq int
	Per  time.Duration
}

func (cp RateLimiter) String() string {
	return fmt.Sprintf("Constant{%d hits/%s}", cp.Freq, cp.Per)
}

func (cp RateLimiter) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	switch {
	case cp.Per == 0 || cp.Freq == 0:
		return 0, false // Zero value = infinite rate
	case cp.Per < 0 || cp.Freq < 0:
		return 0, true
	}

	expectedHits := uint64(cp.Freq) * uint64(elapsed/cp.Per)
	if hits < expectedHits {
		// Running behind, send next hit immediately.
		return 0, false
	}
	interval := uint64(cp.Per.Nanoseconds() / int64(cp.Freq))
	if math.MaxInt64/interval < hits {
		// We would overflow delta if we continued, so stop the attack.
		return 0, true
	}
	delta := time.Duration((hits + 1) * interval)
	// Zero or negative durations cause time.Sleep to return immediately.
	return delta - elapsed, false
}

func (cp RateLimiter) hitsPerNs() float64 {
	return float64(cp.Freq) / float64(cp.Per)
}
