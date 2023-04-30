package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb        *redis.Client
	expiration time.Duration // Zero expiration means the code has no expiration time.
}

func NewRedisStore(options *redis.Options, expiration time.Duration) Store {
	return &RedisStore{
		rdb:        redis.NewClient(options),
		expiration: expiration,
	}
}

func (s *RedisStore) SetIfNotExists(ctx context.Context, code string) (bool, error) {
	set, err := s.rdb.SetNX(ctx, code, nil, s.expiration).Result()
	if err != nil {
		return false, err
	}
	return set, nil
}

func (s *RedisStore) Delete(ctx context.Context, code string) error {
	return s.rdb.Del(ctx, code).Err()
}
