package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	client *redis.Client
}

func NewRedisConnection(addr, password string, db int) (*RedisConnection, error) {
	clientRedis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisConnection{client: clientRedis}, nil
}

func (r *RedisConnection) Get(ctx context.Context, key string) (int64, error) {
	return r.client.Get(ctx, key).Int64()
}

func (r *RedisConnection) Set(ctx context.Context, key, value string, expirationTime time.Duration) error {
	return r.client.Set(ctx, key, value, expirationTime).Err()
}

func (r *RedisConnection) Incr(ctx context.Context, key string) (int64, error) {
	value, err := r.client.Incr(ctx, key).Uint64()
	return int64(value), err
}
