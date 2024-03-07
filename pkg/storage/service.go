package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/pkg/cache"
)

type ServiceInterface interface {
	StoreEventData(events []model.Data) error
}

type StoreDataAPI struct {
	RdbClient *cache.Redis
}

type EventTypeSummary struct {
	CameraID  string `json:"cameraId"`
	Timestamp string `json:"timestamp"`
}

type CameraSummary struct {
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
}

// StoreEventData stores the event Data in a formatted way in redis with a combination of cameraID, timestmap and eventType
func (sd StoreDataAPI) StoreEventData(events []model.Data) error {
	ctx := context.Background()

	placeHolderValue := "_"
	for _, event := range events {
		key := event.Info.Event.CameraID + placeHolderValue + event.Info.Event.Timestamp + placeHolderValue + event.Info.Event.EventType
		err := sd.RdbClient.Add(ctx, key, []byte{}, time.Hour)
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

func GetEventCountSummaryByEventType(eventType string, rdbClient *cache.Redis) (map[string]int64, error) {
	ctx := context.Background()
	keys, err := rdbClient.Scan(ctx, "*_"+eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	eventCounts := make(map[string]int64)
	for _, key := range keys {
		parts := strings.Split(key, "_")
		if len(parts) != 3 {
			continue // Skip if the key format is invalid
		}
		cameraId := parts[0]

		eventCounts[cameraId]++
	}

	return eventCounts, nil
}

func GetEventSummaryByCameraID(cameraId string, rdbClient *cache.Redis) ([]CameraSummary, error) {
	ctx := context.Background()
	keys, err := rdbClient.Scan(ctx, cameraId+"_*")
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	var cameraSummary []CameraSummary
	for _, key := range keys {
		parts := strings.Split(key, "_")
		if len(parts) != 3 {
			continue // Skip if the key format is invalid
		}
		eventType := parts[2]
		timestamp := parts[1]

		event := CameraSummary{EventType: eventType, Timestamp: timestamp}
		cameraSummary = append(cameraSummary, event)
	}

	return cameraSummary, nil
}

func GetEventSummaryByEventType(eventType string, rdbClient *cache.Redis) ([]EventTypeSummary, error) {
	ctx := context.Background()
	keys, err := rdbClient.Scan(ctx, "*_"+eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	var eventSummary []EventTypeSummary
	for _, key := range keys {
		parts := strings.Split(key, "_")
		if len(parts) != 3 {
			continue // Skip if the key format is invalid
		}
		cameraId := parts[0]
		timestamp := parts[1]

		event := EventTypeSummary{CameraID: cameraId, Timestamp: timestamp}
		eventSummary = append(eventSummary, event)
	}

	return eventSummary, nil
}

func GetTotalEventCount(rdbClient *cache.Redis) (int, error) {
	ctx := context.Background()
	keys, err := rdbClient.Scan(ctx, "")
	if err != nil {
		return 0, fmt.Errorf("Error fetching keys: %v", err)
	}

	// Count the number of keys.
	return len(keys), nil
}
