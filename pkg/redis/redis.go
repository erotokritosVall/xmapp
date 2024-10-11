package redis

import (
	"context"
	"fmt"
	"time"

	rd "github.com/redis/go-redis/v9"
)

type Configuration struct {
	Uri string `envconfig:"REDIS_URI"`
}

type Redis interface {
	GetString(ctx context.Context, key string) (*string, error)
	SetString(ctx context.Context, key, value string, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}

type redis struct {
	client *rd.Client
}

func New(ctx context.Context, cfg *Configuration) (Redis, error) {
	opt, err := rd.ParseURL(cfg.Uri)
	if err != nil {
		return nil, fmt.Errorf("redis failed to parse url: %+v", err)
	}

	rd := rd.NewClient(opt)

	if err := rd.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis failed to ping: %+v", err)
	}

	return &redis{
		client: rd,
	}, nil
}

func (r *redis) GetString(ctx context.Context, key string) (*string, error) {
	v, err := r.client.Get(ctx, key).Result()
	if err == rd.Nil {
		return nil, nil
	}

	return &v, err
}

func (r *redis) SetString(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redis) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}
