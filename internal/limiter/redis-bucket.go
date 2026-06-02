package limiter

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisBucket struct {
	Capacity   float64
	ip         string
	rdb        *redis.Client
	refillRate float64
}

func NewRedisBucket(connection *redis.Client, capacity float64, refillRate float64, ip string) *redisBucket {
	return &redisBucket{capacity, ip, connection, refillRate}
}

func (b *redisBucket) Allow() bool {

	ctx := context.Background()

	keys := []string{
		b.ip + ":tokens",
		b.ip + ":last_refill",
	}

	result, err := b.rdb.Eval(ctx, luaScript, keys, b.Capacity, b.refillRate).Int()
fmt.Println("ip is", b.ip)
	if err != nil {
		return false
	}
	fmt.Println("result is ",result)

return result==1

}

var luaScript = `

-- get from redis

local time = redis.call('TIME')
local now = tonumber(time[1])+tonumber(time[2]/1000000)


local tokens = tonumber(redis.call('GET', KEYS[1]))
local last_refill = tonumber(redis.call('GET', KEYS[2]))

local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])

if tokens == false or last_refill == false or last_refill == nil or tokens  ==  nil then
 tokens= capacity
last_refill=now
end

--refill logic


local elapsed=now-last_refill
local current_tokens= elapsed*refill_rate
tokens=math.min(current_tokens+tokens,capacity)
last_refill=now


redis.call('SET',KEYS[1],tokens)
redis.call('SET',KEYS[2],last_refill)

if tokens>=1 then 
tokens=tokens-1
redis.call('SET',KEYS[1],tokens)
return 1
else 
return 0
end

`
