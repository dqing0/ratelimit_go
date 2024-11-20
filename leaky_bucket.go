package ratelimit_go

import (
	"sync/atomic"
	"time"
)

type leakyBucketRateLimiter struct {
	limit        int32
	perReqTime   time.Duration
	prevTakeTime int64
}

func NewLeakyBucket(limit int, opts ...Option) RateLimiter {
	l := &leakyBucketRateLimiter{
		limit:        int32(limit),
		prevTakeTime: 0,
	}
	l.calcPeriodTime()
	for _, opt := range opts {
		opt.apply(l)
	}
	return l
}

func (l *leakyBucketRateLimiter) Take() bool {
	var (
		now         int64
		newNextTime int64
	)
	for {
		now = time.Now().UnixNano()
		prevTime := atomic.LoadInt64(&l.prevTakeTime)

		if prevTime == 0 || (now-prevTime) > int64(l.perReqTime) {
			newNextTime = now
		} else {
			newNextTime = prevTime + int64(l.perReqTime)
		}

		if atomic.CompareAndSwapInt64(&l.prevTakeTime, prevTime, newNextTime) {
			break
		}
	}

	deltaDuration := newNextTime - now
	if deltaDuration > 0 {
		time.Sleep(time.Duration(deltaDuration))
	}
	return true
}

func (l *leakyBucketRateLimiter) calcPeriodTime() {
	periodDuration := time.Second / time.Duration(l.limit)
	l.perReqTime = periodDuration
}
