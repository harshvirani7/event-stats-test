package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/storage"
)

func TotalEventCountByCameraId(c *gin.Context) {
	// Get eventType from query parameter
	cameraId := c.Query("cameraId")
	if cameraId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cameraId query parameter is required"})
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
		"cameraId":          cameraId,
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
