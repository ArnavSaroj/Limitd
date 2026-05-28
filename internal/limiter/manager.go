package limiter

import (
	"fmt"
	"sync"
)

type Manager struct {
	Buckets    map[string]*Bucket
	Mutex      sync.Mutex
	Capacity   float64
	RefillRate float64
}

func NewManager(capacity float64, refillRate float64) *Manager {
	return &Manager{Buckets: make(map[string]*Bucket), Capacity: capacity, RefillRate: refillRate}

}

func (m *Manager) GetBucket(ip string) *Bucket {

	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	bucket, exists := m.Buckets[ip]

	if !exists {
		bucket = newBucket(m.Capacity, m.RefillRate)
		fmt.Print("ip mapping doesnt exists creating new bucket for it")
		m.Buckets[ip] = bucket

	}

	return bucket

}
