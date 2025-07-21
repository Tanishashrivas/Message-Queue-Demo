package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDISHOST"), os.Getenv("REDISPORT")),
	})

	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) PushToQueue(ctx context.Context, queue string, message string) error {
	return r.client.LPush(ctx, queue, message).Err()
}

func (r *RedisClient) ConsumeFromQueue(ctx context.Context, queue string) (string, error) {
	result, err := r.client.BRPop(ctx, 0*time.Second, queue).Result()
	if err != nil || len(result) < 2 {
		return "", err
	}
	return result[1], nil
}
