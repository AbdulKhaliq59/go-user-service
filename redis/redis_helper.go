package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisHelper struct {
	client  *redis.Client
	appName string
}

func NewRedisHelper() *RedisHelper {
	appName := os.Getenv("APPLICATION_NAME")
	if appName == "" {
		appName = "go-transaction-service"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &RedisHelper{
		client:  client,
		appName: appName,
	}
}

func (r *RedisHelper) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if ttl == 0 {
		ttl = 10 * time.Minute // Default 10 minutes
	}

	return r.client.Set(ctx, fmt.Sprintf("%s-%s", r.appName, key), jsonData, ttl).Err()
}

func (r *RedisHelper) Get(key string, dest interface{}) error {
	ctx := context.Background()
	val, err := r.client.Get(ctx, fmt.Sprintf("%s-%s", r.appName, key)).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisHelper) Del(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, fmt.Sprintf("%s-%s", r.appName, key)).Err()
}

func (r *RedisHelper) Unlink(key string) error {
	ctx := context.Background()
	return r.client.Unlink(ctx, fmt.Sprintf("%s-%s", r.appName, key)).Err()
}
