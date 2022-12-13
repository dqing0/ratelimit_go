package ratelimit_go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTokenBucketRateLimiter(t *testing.T) {
	// 10 concurrent and QPS is 1000
	concurrent := 10
	qps := 1000
	limiter := NewRateLimit(qps, WithInterval(100*time.Millisecond))
	var wg sync.WaitGroup
	wg.Add(concurrent)

	var count, stopFlag int32
	startTime := time.Now()
	for i := 0; i < concurrent; i++ {
		go func(tmp int) {
			for atomic.LoadInt32(&stopFlag) == 0 {
				fmt.Printf("concurrent:%v start take\n", tmp)
				if limiter.Take() {
					atomic.AddInt32(&count, 1)
					fmt.Printf("concurrent:%v start take success\n", tmp)
				}
			}
			wg.Done()
		}(i)
	}
	time.AfterFunc(time.Second*2, func() {
		atomic.StoreInt32(&stopFlag, 1)
	})

	wg.Wait()
	spanTime := int32(time.Since(startTime).Seconds())
	actualQPS := count / spanTime
	fmt.Printf("actualQps:%v", actualQPS)
	delta := math.Abs(float64(actualQPS - int32(qps)))
	assert.True(t, delta < float64(qps)*0.02, "rateLimit worked!")
}
