package storage

import (
	"context"
	"fmt"
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
func GetEventCountByCameraID(eventType string) (map[string]int64, error) {
	ctx := context.Background()
	// Use Redis HGetAll command to get all the fields and values in the hash set
	fieldsValues := redisClient.HGetAll(ctx, "event:"+eventType).Val()

	// Initialize a map to store counts for each cameraid
	counts := make(map[string]int64)

	// Loop through the fields and values and calculate the count for each cameraid
	for field, _ := range fieldsValues {
		counts[field]++
	}

	return counts, nil
}

// Function to get total event count by cameraId from Redis
func GetTotalEventCountByCameraID(cameraId string) (int64, error) {
	ctx := context.Background()
	// Use Redis HLen command to get the count of keys in the hash set
	count := redisClient.HLen(ctx, "event:"+cameraId).Val()
	return count, nil
}

// Function to get all data from Redis
func GetAllEventData() (map[string]string, error) {
	ctx := context.Background()
	// Use Redis HGetAll command to get all the fields and values in the hash set
	data := redisClient.HGetAll(ctx, "events").Val()
	return data, nil
}

// update key to be cameraID_timestamp_eventype

// *eventtype - regex for search
