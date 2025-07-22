package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDISHOST"), os.Getenv("REDISPORT"))
	log.Printf("Connecting to Redis at: %s", addr)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis")
	}

	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) PushToQueue(ctx context.Context, queue string, message string) error {
	log.Printf("Pushing message to queue '%s': %s", queue, message)
	err := r.client.LPush(ctx, queue, message).Err()
	if err != nil {
		log.Printf("Error pushing to queue: %v", err)
		return err
	}

	length, _ := r.client.LLen(ctx, queue).Result()
	log.Printf("Queue '%s' length after push: %d", queue, length)
	return nil
}

func (r *RedisClient) ConsumeFromQueue(ctx context.Context, queue string) (string, error) {
	result, err := r.client.BRPop(ctx, 0*time.Second, queue).Result()
	if err != nil || len(result) < 2 {
		return "", err
	}

	message := result[1]
	log.Printf("Consumed message from queue '%s': %s", queue, message)

	length, _ := r.client.LLen(ctx, queue).Result()
	log.Printf("Queue '%s' remaining length: %d", queue, length)

	return message, nil
}
