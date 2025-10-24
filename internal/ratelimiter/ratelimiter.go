package ratelimiter

import "time"

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	ReuestsPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}

func NewFixedWindowLimiter(rpf int, tf time.Duration) Limiter {
	return &FixedWindowLimiter{cfg: cfg}
}
