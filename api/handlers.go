package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
