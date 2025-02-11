package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	ttl    time.Duration
	client *redis.Client
}

func NewRedisCache(host string, port int, password string, db int, ttl int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: client, ttl: time.Duration(ttl) * time.Second}
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(context.TODO(), key).Result()
}

func (r *RedisCache) Set(key string, value string) error {
	return r.client.Set(context.TODO(), key, value, r.ttl).Err()
}
