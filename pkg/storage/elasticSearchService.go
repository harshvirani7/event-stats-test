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
)

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
		NumWorkers:    4,
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

// Function to get total event count by eventType from Elasticsearch
func GetTotalEventCountByTypeES(eventType string, esClient *elasticsearch.Client) (int64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"info.event.eventType.keyword": eventType,
			},
		},
	}

	// Execute the Elasticsearch query
	response, err := executeElasticsearchQuery(query, "event_stats_data", esClient)
	if err != nil {
		return 0, err
	}

	hits := response["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	count := int64(hits)

	return count, nil
}

// Function to get event count by cameraid for a given eventType from Elasticsearch
func GetEventCountByCameraIDES(cameraId string, esClient *elasticsearch.Client) (int64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"info.event.cameraid.keyword": cameraId,
			},
		},
	}

	// Execute the Elasticsearch query
	response, err := executeElasticsearchQuery(query, "event_stats_data", esClient)
	if err != nil {
		return 0, err
	}

	hits := response["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	count := int64(hits)

	return count, nil
}

func GetEventCountSummaryByCameraIDES(cameraId string, esClient *elasticsearch.Client) (map[string]int64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"info.event.cameraid.keyword": cameraId,
			},
		},
		"size": 10000,
	}

	// Execute the Elasticsearch query
	response, err := executeElasticsearchQuery(query, "event_stats_data", esClient)
	if err != nil {
		return nil, err
	}

	// Aggregate event counts by eventType
	eventCounts := make(map[string]int64)
	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		eventType := source["info"].(map[string]interface{})["event"].(map[string]interface{})["eventType"].(string)
		eventCounts[eventType]++
	}

	return eventCounts, nil
}

func GetEventCountSummaryByEventTypeES(eventType string, esClient *elasticsearch.Client) (map[string]int64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"info.event.eventType.keyword": eventType,
			},
		},
		"size": 10000,
	}

	// Execute the Elasticsearch query
	response, err := executeElasticsearchQuery(query, "event_stats_data", esClient)
	if err != nil {
		return nil, err
	}

	// Aggregate event counts by cameraID
	eventCounts := make(map[string]int64)
	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		cameraID := source["info"].(map[string]interface{})["event"].(map[string]interface{})["cameraid"].(string)
		eventCounts[cameraID]++
	}

	return eventCounts, nil
}

func GetEventSummaryByCameraIDES(cameraId string, esClient *elasticsearch.Client) ([]CameraSummary, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"info.event.cameraid.keyword": cameraId,
			},
		},
	}

	// Execute the Elasticsearch query
	response, err := executeElasticsearchQuery(query, "event_stats_data", esClient)
	if err != nil {
		return nil, err
	}

	// Retrieve hits from the response
	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})

	// Parse hits to extract EventType and Timestamp
	var cameraSummary []CameraSummary
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		eventType := source["info"].(map[string]interface{})["event"].(map[string]interface{})["eventType"].(string)
		timestamp := source["info"].(map[string]interface{})["event"].(map[string]interface{})["timestamp"].(string)

		event := CameraSummary{EventType: eventType, Timestamp: timestamp}
		cameraSummary = append(cameraSummary, event)
	}

	return cameraSummary, nil
}

func GetEventSummaryByEventTypeES(eventType string, esClient *elasticsearch.Client) ([]EventTypeSummary, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"info.event.eventType.keyword": eventType,
			},
		},
	}

	// Execute the Elasticsearch query
	response, err := executeElasticsearchQuery(query, "event_stats_data", esClient)
	if err != nil {
		return nil, err
	}

	// Retrieve hits from the response
	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})

	// Parse hits to extract CameraID and Timestamp
	var eventSummary []EventTypeSummary
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		cameraID := source["info"].(map[string]interface{})["event"].(map[string]interface{})["cameraid"].(string)
		timestamp := source["info"].(map[string]interface{})["event"].(map[string]interface{})["timestamp"].(string)

		event := EventTypeSummary{CameraID: cameraID, Timestamp: timestamp}
		eventSummary = append(eventSummary, event)
	}

	return eventSummary, nil
}

func GetTotalEventCountES(esClient *elasticsearch.Client) (int, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	// Execute the Elasticsearch query
	response, err := executeElasticsearchQuery(query, "event_stats_data", esClient)
	if err != nil {
		return 0, err
	}

	// Retrieve the total count from the response
	hits := response["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	count := int(hits)

	return count, nil
}

func executeElasticsearchQuery(query map[string]interface{}, indexName string, esClient *elasticsearch.Client) (map[string]interface{}, error) {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to encode query: %v", err)
	}

	// Perform the search request
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(string(queryBytes))),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search request: %v", err)
	}
	defer res.Body.Close()

	// Check if the search request was successful
	if res.IsError() {
		return nil, fmt.Errorf("search request failed: %s", res.Status())
	}

	// Parse the search response
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %v", err)
	}

	return response, nil
}
