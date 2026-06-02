package limiter

import (

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	//Buckets    map[string]*Bucket--earlier before introducing redis we used this 
    //Mutex      sync.Mutex same as above can be removed cos we are guaranting atomicity through lua
	rdb *redis.Client
	Capacity   float64
	RefillRate float64
}

func NewManager(rdb*redis.Client,capacity float64,refillRate float64)*Manager{
	return &Manager{rdb:rdb,Capacity: capacity,RefillRate: refillRate}
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

return NewRedisBucket(m.rdb,m.Capacity,m.RefillRate,ip)


}
