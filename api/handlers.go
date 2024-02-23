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

	// Return total event count and eventType in the response
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

func EventCountByCameraID(c *gin.Context) {
	// Get eventType from query parameter
	eventType := c.Query("eventType")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})
		return
	}

	// Get total event count for each cameraid for the given eventType
	counts, err := storage.GetEventCountByCameraID(eventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return counts in the response
	c.JSON(http.StatusOK, counts)
}

func TotalEventCountByCameraID(c *gin.Context) {
	// Get cameraId from query parameter
	cameraId := c.Query("cameraId")
	if cameraId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cameraId query parameter is required"})
		return
	}

	// Get total event count for the given cameraId
	count, err := storage.GetTotalEventCountByCameraID(cameraId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return total event count and cameraId in the response
	response := gin.H{
		"cameraId":          cameraId,
		"total_event_count": count,
	}
	c.JSON(http.StatusOK, response)
}

func GetAllEventData(c *gin.Context) {
	// Get all data from Redis
	data, err := storage.GetAllEventData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return data in a nicely formatted response
	c.JSON(http.StatusOK, data)
}
