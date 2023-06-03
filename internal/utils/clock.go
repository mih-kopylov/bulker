package utils

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct {
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

type FixedClock struct {
	nowTime time.Time
}

func (c *FixedClock) Now() time.Time {
	return c.nowTime
}
