package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client with the provided configuration.
func NewRedisClient(config Config) (*redis.Client, error) {
	// Create a new Redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: "",
		DB:       0,
	})

	// Test the connection to the Redis server.
	err := client.Ping(client.Context()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	// Create a new Redis client instance and return it.
	return client, nil
}
