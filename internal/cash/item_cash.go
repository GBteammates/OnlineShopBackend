package cash

import (
	"OnlineShopBackend/internal/handlers"
	"OnlineShopBackend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var _ Cash = &RedisCash{}

type RedisCash struct {
	*redis.Client
	TTL    time.Duration
	logger *zap.Logger
}

type results struct {
	Responses []handlers.Item
}

// NewRedisCash initialize redis client
func NewRedisCash(host, port string, ttl time.Duration, logger *zap.Logger) (*RedisCash, error) {
	logger.Debug("Enter in NewRedisCash()")
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
	cashTTL := ttl * time.Minute
	c := &RedisCash{
		Client: client,
		TTL:    cashTTL,
		logger: logger,
	}
	return c, nil
}

// Shutdown close redis client
func (cash *RedisCash) Close() error {
	cash.logger.Debug("Enter in RedisCash Close()")
	return cash.Client.Close()
}

// CheckCash checks for data in the cache
func (cash *RedisCash) CheckCash(key string) bool {
	cash.logger.Debug("Enter in cash CheckCash()")
	item, err := cash.GetCash(key)
	if err != nil {
		cash.logger.Error(fmt.Sprintf("redis: get record %q: %v", key, err))
		return false
	}

	if item != nil {
		cash.logger.Debug(fmt.Sprintf("Key %q in cash found success", key))
		return true
	}
	cash.logger.Debug(fmt.Sprintf("Redis: get record %q not exist", key))
	return false
}

// CreateCash add data in the cash
func (cash *RedisCash) CreateCash(ctx context.Context, res chan models.Item, key string) error {
	cash.logger.Debug("Enter in cash CreateCash()")
	in := results{}
	for resItem := range res {
		in.Responses = append(in.Responses, handlers.Item{
			Id:          resItem.Id.String(),
			Title:       resItem.Title,
			Description: resItem.Description,
			Price:       resItem.Price,
			Category:    resItem.Category.String(),
			Vendor:      resItem.Vendor,
		})
	}

	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal unknown item: %w", err)
	}

	err = cash.Set(ctx, key, data, cash.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: set key %q: %w", key, err)
	}
	return nil
}

// GetCash retrieves data from the cache
func (cash *RedisCash) GetCash(key string) ([]handlers.Item, error) {
	cash.logger.Debug("Enter in cash GetCash()")
	res := results{}
	data, err := cash.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &res)
	return res.Responses, nil
}
