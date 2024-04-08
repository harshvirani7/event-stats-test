package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/pkg/cache"
)

type ServiceInterface interface {
	StoreEventData(events []model.Data) error
}

type StoreDataAPI struct {
	RdbClient *cache.Redis
	EsClient  *elasticsearch.Client
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

func (sd StoreDataAPI) StoreInElasticsearch(data []model.Data) error {

	var requests []esutil.BulkIndexerItem

	for _, d := range data {
		jsonData, err := json.Marshal(d)
		if err != nil {
			return err
		}

		request := esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: d.Unique,
			Body:       bytes.NewReader(jsonData),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
				fmt.Printf("Failed to index document %s: %s\n", item.DocumentID, err)
			},
		}

		requests = append(requests, request)
	}

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        sd.EsClient,
		Index:         "event_stats_data",
		NumWorkers:    4, // Adjust based on your system resources
		FlushBytes:    10e6,
		FlushInterval: 30 * time.Second,
		OnError:       func(ctx context.Context, err error) { fmt.Printf("Bulk indexer error: %s\n", err) },
	})
	if err != nil {
		return err
	}

	for _, req := range requests {
		if err := bulkIndexer.Add(context.Background(), req); err != nil {
			return err
		}
	}

	if err := bulkIndexer.Close(context.Background()); err != nil {
		return err
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
