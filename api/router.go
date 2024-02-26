package api

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine) {
	// router.POST("/storeEventData", StoreEventData)
	// router.GET("/totalEventCountByEventType", TotalEventCountByType)
	// router.GET("/totalEventCountByCameraId", TotalEventCountByCameraId)
	router.GET("/eventCountSummaryByCameraId", EventCountSummaryByCameraId)
}

// create interface - method to store data
// implement method for pulsar structure\
// interface for stote
// similarly for other handlers

// move to main.go
