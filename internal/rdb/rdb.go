package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
)


type RedisClient struct {
	Rdb *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	return &RedisClient{
		Rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.Rdb.Ping(ctx).Err()
}