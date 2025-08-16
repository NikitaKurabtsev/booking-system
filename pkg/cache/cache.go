package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const TTL = 5 * time.Minute

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type redisCache struct {
	client *redis.Client
}

func NewCache(addr string) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisCache{client: client}, nil
}

func (rc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rc.client.Set(ctx, key, value, ttl).Err()
}

func (rc *redisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

func (rc *redisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

func (rc *redisCache) Close() error {
	return rc.client.Close()
}
