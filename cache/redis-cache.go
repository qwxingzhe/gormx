package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	rdb *redis.Client
}

func InitRedisCache(opt *redis.Options) Interface {
	return &redisCache{
		rdb: redis.NewClient(opt),
	}
}

// SetCache 新建缓存
func (c *redisCache) SetCache(key string, value string) error {
	return c.rdb.Set(context.Background(), key, value, time.Hour*24).Err()
}

// GetCache 查询缓存
func (c *redisCache) GetCache(key string) string {
	return c.rdb.Get(context.Background(), key).Val()
}

// DelCache 删除缓存
func (c *redisCache) DelCache(key string) error {
	return c.rdb.Del(context.Background(), key).Err()
}
