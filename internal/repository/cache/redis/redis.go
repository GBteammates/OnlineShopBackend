package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCache struct {
	*redis.Client
	TTL    time.Duration
	logger *zap.Logger
}

// NewRedisCache initialize redis client
func NewRedisCache(host, port string, ttl time.Duration, logger *zap.Logger) (*RedisCache, error) {
	logger.Sugar().Debugf("Enter in NewRedisCache() with args: host: %s, port: %s, ttl: %v, logger", host, port, ttl)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}
	logger.Debug("Redis Client ping success")
	cacheTTL := ttl * time.Hour
	c := &RedisCache{
		Client: client,
		TTL:    cacheTTL,
		logger: logger,
	}
	return c, nil
}

// ShutDown is func for graceful shutdown redis connection
func (cache *RedisCache) Shutdown(timeout int) error {
	cache.logger.Sugar().Debugf("Enter in cache ShutDown() with args: timeout: %d", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	status := cache.Client.Shutdown(ctx)
	result, err := status.Result()
	if err != nil {
		return err
	}
	cache.logger.Info(result)
	return nil
}
