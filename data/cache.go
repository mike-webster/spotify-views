package data

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
}

type TestCache struct {
	setError error
	getError error
}

type LiveCache struct {
	cache *redis.Client
}

func (c *LiveCache) Set(ctx context.Context, key, value string) error {
	_, err := c.cache.Set(ctx, key, value, 2*time.Hour).Result()
	return err
}

func (c *LiveCache) Get(ctx context.Context, key string) (string, error) {
	return c.cache.Get(ctx, key).Result()
}

func (c *TestCache) Set(ctx context.Context, key, value string) error {
	if c.setError != nil {
		return c.setError
	}

	return nil
}

func (c *TestCache) Get(ctx context.Context, key string) (string, error) {
	if c.getError != nil {
		return "", c.getError
	}

	return "", nil
}

func GetLiveCache(ctx context.Context) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &LiveCache{cache: rdb}, nil
}
