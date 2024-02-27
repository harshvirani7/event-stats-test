package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/pkg/cache"
)

var redisClient *redis.Client

func InitRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
}

// StoreEventData stores the event Data in a formatted way in redis with a combination of cameraID, timestmap and eventType
func StoreEventData(events []model.Data, rdbClient *cache.Redis) error {
	ctx := context.Background()

	placeHolderValue := "_"
	for _, event := range events {
		key := event.Info.Event.CameraID + placeHolderValue + event.Info.Event.Timestamp + placeHolderValue + event.Info.Event.EventType
		err := rdbClient.Add(ctx, key, []byte{}, time.Hour)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to get total event count by eventType from Redis
func GetTotalEventCountByType(eventType string, rdbClient *cache.Redis) (int64, error) {
	ctx := context.Background()
	keys, err := rdbClient.Scan(ctx, "*_"+eventType)
	if err != nil {
		return 0, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	count := int64(len(keys))
	return count, nil
}

// Function to get event count by cameraid for a given eventType from Redis
func GetEventCountByCameraID(cameraId string, rdbClient *cache.Redis) (int64, error) {
	ctx := context.Background()
	keys, err := rdbClient.Scan(ctx, cameraId+"_*")
	if err != nil {
		return 0, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	count := int64(len(keys))
	return count, nil
}

func GetEventCountSummaryByCameraID(cameraId string, rdbClient *cache.Redis) (map[string]int64, error) {
	ctx := context.Background()
	keys, err := rdbClient.Scan(ctx, cameraId+"_*")
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

	return eventCounts, nil
}
