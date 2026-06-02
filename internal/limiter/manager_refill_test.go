package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/arnavsaroj/goratelimiter/internal/limiter"
	"github.com/arnavsaroj/goratelimiter/internal/store"
)

func TestRefill(t *testing.T) {
	rdb := store.NewRedisConnection()
	ctx := context.Background()

	rdb.Del(ctx, "test-user:tokens")
	rdb.Del(ctx, "test-user:last_refill")

	manager := limiter.NewManager(rdb, 10, 1) // capacity=10, refill=1 token/sec

	bucket := manager.GetBucket("test-user")

	// Consume all 10 tokens
	for i := 0; i < 10; i++ {
		if !bucket.Allow() {
			t.Fatalf("expected request %d to be allowed", i)
		}
	}

	// Next one should fail
	if bucket.Allow() {
		t.Fatal("expected bucket to be empty")
	}

	t.Log("Bucket exhausted")

	time.Sleep(5 * time.Second)

	allowed := 0

	for i := 0; i < 10; i++ {
		if bucket.Allow() {
			allowed++
		}
	}

	t.Logf("Allowed after refill: %d", allowed)

	if allowed < 4 || allowed > 6 {
		t.Fatalf("expected around 5 tokens after 5 seconds, got %d", allowed)
	}
}