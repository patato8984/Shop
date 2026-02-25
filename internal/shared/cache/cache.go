package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) bool
	Get(ctx context.Context, key string, dest interface{}) bool
	Del(ctx context.Context, key string) bool
	PipelineDelCart(ctx context.Context, key []string) error
}
type RedisCache struct {
	log    *zap.Logger
	client *redis.Client
}

func NewRedisCache(client *redis.Client, log *zap.Logger) Cache {
	return &RedisCache{client: client, log: log}
}
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) bool {
	data, err := json.Marshal(value)
	if err != nil {
		c.log.Error("cache",
			zap.Error(err),
		)
		return false
	}
	_, err = c.client.Set(ctx, key, data, ttl).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		c.log.Error("cache",
			zap.Error(err),
		)
		return false
	}
	return true
}
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) bool {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		c.log.Error("cache",
			zap.Error(err),
		)
		return false
	}
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return false
	}
	return true
}
func (c *RedisCache) Del(ctx context.Context, key string) bool {
	_, err := c.client.Del(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		c.log.Error("cache",
			zap.Error(err),
		)
		return false
	}
	return true
}
func (c *RedisCache) PipelineDelCart(ctx context.Context, key []string) error {
	pipe := c.client.Pipeline()
	for _, idCart := range key {
		pipe.Del(ctx, "cart:"+idCart)
	}
	_, err := pipe.Exec(ctx)
	return err
}
func GetOrSet[T any](
	ctx context.Context,
	c Cache,
	key string,
	ttl time.Duration,
	dest *T,
	dbFunc func() (*T, error),
) (*T, error) {
	ok := c.Get(ctx, key, dest)
	if ok {
		return dest, nil
	}
	result, err := dbFunc()
	if err != nil {
		return nil, err
	}
	_ = c.Set(ctx, key, result, ttl)
	return result, nil
}
