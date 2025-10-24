package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/config"
)

// RedisCache wraps the Redis client
type RedisCache struct {
	Client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis connection
func NewRedisCache(cfg *config.Config) (*RedisCache, error) {
	// Parse Redis URL or build from components
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		// Fallback to manual configuration
		opt = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       0, // default DB
		}
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis cache")

	return &RedisCache{
		Client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis with expiration
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(r.ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis
func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.Client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

// Delete removes a key from Redis
func (r *RedisCache) Delete(key string) error {
	return r.Client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisCache) Exists(key string) (bool, error) {
	val, err := r.Client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// SetJSON stores a JSON-serializable value in Redis
func (r *RedisCache) SetJSON(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(r.ctx, key, value, expiration).Err()
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.Client.Close()
}

// HealthCheck verifies the Redis connection is healthy
func (r *RedisCache) HealthCheck() error {
	return r.Client.Ping(r.ctx).Err()
}
