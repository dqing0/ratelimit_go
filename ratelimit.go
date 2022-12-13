package ratelimit_go

type Option struct {
	apply func(RateLimiter)
}

type RateLimiter interface {
	Take() bool
}

func NewRateLimit(limit int, opts ...Option) RateLimiter {
	return NewRateLimitWithAlgorithm(limit, "token_bucket", opts...)
}

func NewRateLimitWithAlgorithm(limit int, algorithm string, opts ...Option) RateLimiter {
	switch algorithm {
	case "token_bucket":
		return NewTokenBucket(limit, opts...)
	case "leaky_bucket":
		return NewLeakyBucket(limit, opts...)
	}
	return nil
}
