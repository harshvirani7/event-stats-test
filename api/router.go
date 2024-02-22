package api

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine) {
	router.GET("/totalEventCount", TotalEventCount)
	router.POST("/storeEventData", StoreEventData)
	// Define other routes...
}
