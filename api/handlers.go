package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/storage"
)

func TotalEventCountByType(c *gin.Context) {
	// Get event type from query parameter
	eventType := c.Query("eventType")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})
		return
	}

	// Get total event count for the given event type
	count, err := storage.GetTotalEventCountByType(eventType)
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

func StoreEventData(c *gin.Context) {
	var data []model.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store events into Redis
	if err := storage.StoreEventData(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event data stored successfully"})
}

func TotalEventCountByCameraId(c *gin.Context) {
	// Get eventType from query parameter
	cameraId := c.Query("cameraId")
	if cameraId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})
		return
	}

	// Get total event count for each cameraid for the given eventType
	count, err := storage.GetEventCountByCameraID(cameraId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return total count and eventType in the response
	response := gin.H{
		"eventType":         cameraId,
		"total_event_count": count,
	}
	c.JSON(http.StatusOK, response)
}

func EventCountSummaryByCameraId(c *gin.Context) {
	// Get cameraId from query parameter
	cameraId := c.Query("cameraId")
	if cameraId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cameraId query parameter is required"})
		return
	}

	// Get event counts for the given cameraId from Redis
	eventCounts, err := storage.GetEventCountSummaryByCameraID(cameraId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return event counts for the given cameraId
	c.JSON(http.StatusOK, eventCounts)
}
