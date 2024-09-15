package cache

import (
	"context"
	"errors"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

// TODO : Implement MULTI-EXEC for atomic ops.

// RedisCache struct implementing ICache interface
type RedisCache[T any] struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration // Time-to-live for cache items
	log    logger.ILogger
}

// NewRedisCache returns a new RedisCache with provided config
func NewRedisCache[T any](config *Config, logger logger.ILogger) *RedisCache[T] {
	return &RedisCache[T]{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		}),
		ctx: context.Background(),
		ttl: config.TTL,
		log: logger,
	}
}

// Get method retrieves value by key from Redis
func (r *RedisCache[T]) Get(key interface{}) (T, bool) {
	var result T

	// Convert key to string (assuming Redis keys are stored as strings)
	keyStr := key.(string)

	// Try to get the value from Redis
	val, err := r.client.Get(r.ctx, keyStr).Result()
	if errors.Is(err, redis.Nil) {
		// Key does not exist
		return result, false
	} else if err != nil {
		// Handle other errors (optional: log or handle differently)
		return result, false
	}

	// Assuming the value can be unmarshaled into T type
	// Example here for simplicity: you might need to convert it to a specific type.
	result = any(val).(T)
	return result, true
}

// Set method stores key-value pair in Redis with TTL
func (r *RedisCache[T]) Set(key interface{}, val T) bool {
	// Convert key to string (assuming Redis keys are stored as strings)
	keyStr := key.(string)

	// Set the value with TTL in Redis
	err := r.client.Set(r.ctx, keyStr, val, r.ttl).Err()
	if err != nil {
		// Handle error (optional: log or handle differently)
		r.log.Errorf("redis set error: %v", err)
		return false
	}
	return true
}
