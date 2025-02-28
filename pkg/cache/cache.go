package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
	ctx    context.Context
}

func NewCache(addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Cache{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *Cache) Set(key string, value interface{}) error {
	return c.client.Set(c.ctx, key, value, c.ttl).Err()
}

func (c *Cache) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}
