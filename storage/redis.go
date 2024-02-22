package storage

import "github.com/go-redis/redis/v8"

var redisClient *redis.Client

func InitRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
}

func GetTotalEventCount() (int, error) {
	// Implement function to get total event count from Redis
	return 0, nil
}
