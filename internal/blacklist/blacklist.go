package blacklist

import (
	"avto-crm-api/internal/rdb"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type BlacklistService struct {
	client *rdb.RedisClient
}

func NewBlacklistService(client *rdb.RedisClient) *BlacklistService {
	return &BlacklistService{
		client: client,
	}
}

func (b *BlacklistService) Add(ctx context.Context, token string, ttl time.Duration) error {
	return b.client.Rdb.Set(ctx, "blacklist:"+token, "1", ttl).Err()
} 

func (b *BlacklistService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	err := b.client.Rdb.Get(ctx, "blacklist:"+token).Err()

	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}