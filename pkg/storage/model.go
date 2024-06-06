package storage

import (
	"github.com/elastic/go-elasticsearch/v8"
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
