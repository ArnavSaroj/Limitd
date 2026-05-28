package limiter

import (
	"sync"
	"time"
)

type Bucket struct {
	mu sync.Mutex
	NoOfTokens     float64
	RefillRate     float64
	Capacity       float64
	LastRefillTime time.Time
}

func newBucket(capacity float64, refillrate float64) *Bucket {
	return &Bucket{NoOfTokens: capacity, RefillRate: refillrate, Capacity: capacity, LastRefillTime: time.Now()}
}

func (b *Bucket) refill() {

	now := time.Now()

	if b.NoOfTokens == b.Capacity {
		b.LastRefillTime = now
		return
	}

	elapsed := now.Sub(b.LastRefillTime)
	var tokensToAdd float64 = (elapsed.Seconds()* b.RefillRate)
	b.NoOfTokens = min(b.NoOfTokens+tokensToAdd, b.Capacity)

	b.LastRefillTime = now

}

func (b *Bucket) Allow() bool {
	b.mu.Lock()
defer b.mu.Unlock()
	b.refill()
	if b.NoOfTokens > 0 {
		b.NoOfTokens = b.NoOfTokens - 1
		return true
	}

	return false
}
