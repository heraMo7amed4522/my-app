package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache() *RedisCache {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0 // Default Redis DB

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return &RedisCache{
			client: nil,
			ctx:    ctx,
		}
	}

	log.Println("Successfully connected to Redis")

	return &RedisCache{
		client: rdb,
		ctx:    ctx,
	}
}

// Set stores a value in Redis with expiration
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	err = r.client.Set(r.ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}

// Get retrieves a value from Redis
func (r *RedisCache) Get(key string, dest interface{}) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get cache: %v", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %v", err)
	}

	return nil
}

// Delete removes a key from Redis
func (r *RedisCache) Delete(key string) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %v", err)
	}

	return nil
}

// DeletePattern removes all keys matching a pattern
func (r *RedisCache) DeletePattern(pattern string) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %v", err)
	}

	if len(keys) > 0 {
		err = r.client.Del(r.ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete keys: %v", err)
		}
	}

	return nil
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
