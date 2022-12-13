package ratelimit_go

import (
	"sync/atomic"
	"time"
)

var fixedWindowTime = time.Second

type tokenBucketRateLimiter struct {
	limit    int32
	tokens   int32
	interval time.Duration
	once     int32
}

func NewTokenBucket(limit int, opts ...Option) RateLimiter {
	l := &tokenBucketRateLimiter{
		limit: int32(limit),
	}
	for _, opt := range opts {
		opt.apply(l)
	}
	l.calcOnce()
	go l.createTokens()
	return l
}

func WithInterval(interval time.Duration) Option {
	return Option{apply: func(limiter RateLimiter) {
		l := limiter.(*tokenBucketRateLimiter)
		l.interval = interval
	}}
}

func (l *tokenBucketRateLimiter) createTokens() {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		l._createTokens()
	}
}

func (l *tokenBucketRateLimiter) _createTokens() {
	if atomic.LoadInt32(&l.tokens) > l.limit {
		return
	}

	cur := atomic.LoadInt32(&l.tokens)
	if cur+l.once > l.limit {
		atomic.StoreInt32(&l.tokens, l.limit)
	} else {
		atomic.StoreInt32(&l.tokens, l.once)
	}
}

func (l *tokenBucketRateLimiter) Take() bool {
	if atomic.LoadInt32(&l.tokens) <= 0 {
		return false
	}
	return atomic.AddInt32(&l.tokens, -1) >= 0
}

func (l *tokenBucketRateLimiter) calcOnce() {
	if l.interval > time.Second || l.interval == 0 {
		l.interval = 100 * time.Millisecond
	}
	once := int32(float64(l.limit) / (fixedWindowTime.Seconds() / l.interval.Seconds()))
	if once < 1 {
		once = 1
	}
	l.once = once
	l.tokens = once
}
