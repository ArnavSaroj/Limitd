package limiter

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	localBuckets map[string]*Bucket //--earlier before introducing redis we used this
	// i will reintroduce this  for high availability as temp fallback in case redis goes down
	Mutex sync.RWMutex
	//same as above can be removed cos we are guaranting atomicity through lua
	rdb        *redis.Client
	Capacity   float64
	RefillRate float64

	//redis health checker which pings redis every 2 seconds to see if it can respond and if not we fallback to localbuckets
	RedisHealthy atomic.Bool
	//now that we hav health chcker which pings redis and fall back to local memory lets do circuit breaker,it basically like circuit breaker in house
	//now for eg redis failed for 3 req now 4th req will also run the same code and print redis is unhealthy ->wasted cpu cycles hence i implement circuit breaker which immediately goes after 3 failures

	ConsecutiveFailures atomic.Int64
}

func NewManager(rdb *redis.Client, capacity float64, refillRate float64) *Manager {
	m := &Manager{localBuckets: make(map[string]*Bucket), rdb: rdb, Capacity: capacity, RefillRate: refillRate}

	m.RedisHealthy.Store(true)

	return m
}

func (m *Manager) GetBucket(ip string) *redisBucket {

	// m.Mutex.Lock()--all prev ones are commented out by me
	// defer m.Mutex.Unlock()
	// bucket, exists := m.Buckets[ip]

	// if !exists {
	// 	bucket = newBucket(m.Capacity, m.RefillRate)
	// 	fmt.Print("ip mapping doesnt exists creating new bucket for it")
	// 	// m.Buckets[ip] = bucket

	// }

	// return bucket

	return NewRedisBucket(m.rdb, m.Capacity, m.RefillRate, ip)

}

func (m *Manager) GetLocalBucket(ip string) *Bucket {
	m.Mutex.RLock()

	bucket, exists := m.localBuckets[ip]
	m.Mutex.RUnlock()

	if exists {
		return bucket
	}


	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	bucket, exists = m.localBuckets[ip]
	//checking twice here cos imagine if goroutine a and b both arrive at the same time and find that bucket doesnt exists then a will acquire lock create bucket and then release it then b will acquire lock and will create or ovewrite same bucket and then release it means bucket overwritten,so to prevent that we use this double checking

	if !exists {
		bucket = newBucket(m.Capacity, m.RefillRate)

		m.localBuckets[ip] = bucket
	}

	return bucket
}

func (m *Manager) StartRedisHealthchecker() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				500*time.Millisecond,
			)
			err := m.rdb.Ping(ctx).Err()

			cancel()

			if err != nil {

				
				if m.RedisHealthy.Load() {
					fmt.Println("Redis became unhealthy")
				}
					m.RedisHealthy.Store(false)

			} else {
				if !m.RedisHealthy.Load() {
					fmt.Println("Redis recovered")
				}
				m.RedisHealthy.Store(true)
				m.ConsecutiveFailures.Store(0)
			}
		}

	}()
}
