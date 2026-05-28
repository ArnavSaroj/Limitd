package limiter

import (
	"fmt"
	"sync"
	"testing"

)

func TestConcurrentRequests(t *testing.T) {
	manager := NewManager(10, 10.0/60.0) // 10 tokens, 10 per minute

	var wg sync.WaitGroup
	allowed := 0
	denied := 0
	var mu sync.Mutex

	for i := 0; i < 20; i++ { // 20 requests, only 10 should pass
		wg.Add(1)
		go func() {
			defer wg.Done()
			bucket := manager.GetBucket("192.168.1.1") // same IP
			if bucket.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			} else {
				mu.Lock()
				denied++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Allowed: %d, Denied: %d\n", allowed, denied)

	if allowed > 10 {
		t.Errorf("Expected max 10 allowed, got %d", allowed)
	}
}