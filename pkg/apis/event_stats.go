package apis

import (
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/pkg/cache"
	"github.com/harshvirani7/event-stats-test/pkg/config"
	"github.com/harshvirani7/event-stats-test/pkg/monitor"
	"github.com/harshvirani7/event-stats-test/pkg/storage"
	"github.com/harshvirani7/event-stats-test/utils"
	"go.uber.org/zap"
)

type EventStats struct {
	Logger    *zap.SugaredLogger
	RdbClient *cache.Redis
	Cfg       config.Config
	Metrics   *monitor.Metrics
	EsClient  *elasticsearch.Client
}

func (es EventStats) StoreEventData() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var data []model.Data
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sd := storage.StoreDataAPI{
			RdbClient: es.RdbClient,
			EsClient:  es.EsClient,
		}
		// // Store events into Redis
		// if err := sd.StoreEventData(data); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }

		// Store events into Elasticsearch
		if err := sd.StoreInElasticsearch(data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Event data stored successfully"})

		// count, err := storage.GetTotalEventCount(es.RdbClient)
		count, err := storage.GetTotalEventCountES(es.EsClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
		es.Metrics.EventCount.Set(float64(count))
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) TotalEventCountByType() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get event type from query parameter
		eventType := c.Query("eventType")
		if eventType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			return
		}

		// Get total event count for the given event type
		// count, err := storage.GetTotalEventCountByType(eventType, es.RdbClient)
		count, err := storage.GetTotalEventCountByTypeES(eventType, es.EsClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		utils.Sleep(200)

		// Return total count and eventType in the response
		response := gin.H{
			"eventType":         eventType,
			"total_event_count": count,
		}
		c.JSON(http.StatusOK, response)
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) TotalEventCountByCameraId() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get eventType from query parameter
		cameraId := c.Query("cameraId")
		if cameraId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cameraId query parameter is required"})

			return
		}

		// Get total event count for each cameraid for the given eventType
		// count, err := storage.GetEventCountByCameraID(cameraId, es.RdbClient)
		count, err := storage.GetEventCountByCameraIDES(cameraId, es.EsClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		utils.Sleep(200)

		// Return total count and eventType in the response
		response := gin.H{
			"cameraId":          cameraId,
			"total_event_count": count,
		}
		c.JSON(http.StatusOK, response)
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) EventCountSummaryByCameraId() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get cameraId from query parameter
		cameraId := c.Query("cameraId")
		if cameraId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cameraId query parameter is required"})

			return
		}

		// Get event counts for the given cameraId from Redis
		// eventCounts, err := storage.GetEventCountSummaryByCameraID(cameraId, es.RdbClient)
		eventCounts, err := storage.GetEventCountSummaryByCameraIDES(cameraId, es.EsClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return each event counts for the given cameraId
		response := gin.H{
			"cameraId":      cameraId,
			"event_summary": eventCounts,
		}
		c.JSON(http.StatusOK, response)
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) EventCountSummaryByEventType() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get eventType from query parameter
		eventType := c.Query("eventType")
		if eventType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			return
		}

		// Get total camera event counts for the given eventType from Redis
		// eventCounts, err := storage.GetEventCountSummaryByEventType(eventType, es.RdbClient)
		eventCounts, err := storage.GetEventCountSummaryByEventTypeES(eventType, es.EsClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return each event counts for the given cameraId
		response := gin.H{
			"eventType":      eventType,
			"camera_summary": eventCounts,
		}
		c.JSON(http.StatusOK, response)
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) SummaryByCameraId() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get eventType from query parameter
		cameraId := c.Query("cameraId")
		if cameraId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			return
		}

		// Get event summary for the given eventType from Redis
		// eventSummary, err := storage.GetEventSummaryByCameraID(cameraId, es.RdbClient)
		eventSummary, err := storage.GetEventSummaryByCameraIDES(cameraId, es.EsClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		// Return event summary for the given eventType
		c.JSON(http.StatusOK, gin.H{"eventType": cameraId, "event_summary": eventSummary})
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) SummaryByEventType() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get eventType from query parameter
		eventType := c.Query("eventType")
		if eventType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			return
		}

		// Get event summary for the given eventType from Redis
		// eventSummary, err := storage.GetEventSummaryByEventType(eventType, es.RdbClient)
		eventSummary, err := storage.GetEventSummaryByEventTypeES(eventType, es.EsClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		// Return event summary for the given eventType
		c.JSON(http.StatusOK, gin.H{"eventType": eventType, "event_summary": eventSummary})
	}
	return gin.HandlerFunc(fn)
}

func RemovePathParam(c *gin.Context) string {
	var newPath string
	for _, str := range strings.Split(c.Request.URL.Path, "/") {
		found := false
		for _, paramValue := range c.Params {
			if str == paramValue.Value {
				newPath = newPath + "*/"
				found = true
				break
			}
		}
		if !found {
			newPath = newPath + str + "/"
		}
	}
	return newPath
}
