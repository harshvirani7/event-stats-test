package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/pkg/cache"
	"github.com/harshvirani7/event-stats-test/pkg/config"
	"github.com/harshvirani7/event-stats-test/pkg/storage"
	"go.uber.org/zap"
)

type EventStats struct {
	Logger    *zap.SugaredLogger
	RdbClient *cache.Redis
	Cfg       config.Config
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
		}
		// Store events into Redis
		// change this to call from a struct
		if err := sd.StoreEventData(data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Event data stored successfully"})
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
		count, err := storage.GetTotalEventCountByType(eventType, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
		count, err := storage.GetEventCountByCameraID(cameraId, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
		eventCounts, err := storage.GetEventCountSummaryByCameraID(cameraId, es.RdbClient)
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
		eventCounts, err := storage.GetEventCountSummaryByEventType(eventType, es.RdbClient)
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
