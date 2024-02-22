package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/storage"
)

func TotalEventCount(c *gin.Context) {
	count, err := storage.GetTotalEventCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_event_count": count})
}

func StoreEventData(c *gin.Context) {
	var data model.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store eventType, cameraid, and timestamp into Redis
	err := storage.StoreEventData(data.Info.Event.EventType, data.Info.Event.CameraID, data.Info.Event.Timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event data stored successfully"})
}
