package ratelimiter

import "time"

type FixedWindowLimiter struct {
	ReuestsPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}

func (limiter *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	//
}
