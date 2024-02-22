package storage

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/harshvirani7/event-stats-test/model"
)

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

func StoreEventData(events []model.Data) error {
	ctx := context.Background()
	for _, event := range events {
		key := "event:" + event.Info.Event.EventType
		err := redisClient.HSet(ctx, key, event.Info.Event.CameraID, event.Info.Event.Timestamp).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
