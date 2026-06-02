package limiter_test

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/arnavsaroj/goratelimiter/internal/limiter"
	"github.com/arnavsaroj/goratelimiter/internal/store"
)

func TestMultipleManagers(t *testing.T) {
	rdb := store.NewRedisConnection()
	ctx := context.Background()

	rdb.Del(ctx, "distributed-test:tokens")
	rdb.Del(ctx, "distributed-test:last_refill")

	manager1 := limiter.NewManager(rdb, 10, 0)
	manager2 := limiter.NewManager(rdb, 10, 0)
	manager3 := limiter.NewManager(rdb, 10, 0)

	var allowed int64
	var denied int64

	var wg sync.WaitGroup

	totalRequests := 1000

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			var allowedReq bool

			switch rand.Intn(3) {
			case 0:
				allowedReq = manager1.GetBucket("distributed-test").Allow()
			case 1:
				allowedReq = manager2.GetBucket("distributed-test").Allow()
			case 2:
				allowedReq = manager3.GetBucket("distributed-test").Allow()
			}

			if allowedReq {
				atomic.AddInt64(&allowed, 1)
			} else {
				atomic.AddInt64(&denied, 1)
			}
		}()
	}

	wg.Wait()

	t.Logf("Allowed: %d", allowed)
	t.Logf("Denied: %d", denied)

	if allowed != 10 {
		t.Fatalf("expected exactly 10 allowed requests, got %d", allowed)
	}

	if denied != int64(totalRequests)-10 {
		t.Fatalf(
			"expected %d denied requests, got %d",
			totalRequests-10,
			denied,
		)
	}
}
