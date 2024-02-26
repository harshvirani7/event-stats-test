package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

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

// Function to get total event count by eventType from Redis
func GetTotalEventCountByType(eventType string) (int64, error) {
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, "*_"+eventType).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	count := int64(len(keys))
	return count, nil
}

// StoreEventData stores the event Data in a formatted way in redis with a combination of cameraID, timestmap and eventType
func StoreEventData(events []model.Data) error {
	ctx := context.Background()

	placeHolderValue := "_"
	for _, event := range events {
		key := event.Info.Event.CameraID + placeHolderValue + event.Info.Event.Timestamp + placeHolderValue + event.Info.Event.EventType
		err := redisClient.Set(ctx, key, placeHolderValue, time.Hour).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to get event count by cameraid for a given eventType from Redis
func GetEventCountByCameraID(cameraId string) (int64, error) {
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, cameraId+"_*").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	count := int64(len(keys))
	return count, nil
}

func GetEventCountSummaryByCameraID(cameraId string) (map[string]int64, error) {
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, cameraId+"_*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	eventCounts := make(map[string]int64)
	for _, key := range keys {
		parts := strings.Split(key, "_")
		if len(parts) != 3 {
			continue // Skip if the key format is invalid
		}
		eventType := parts[2]

		eventCounts[eventType]++
	}

	fmt.Print(eventCounts)
	return eventCounts, nil
}
