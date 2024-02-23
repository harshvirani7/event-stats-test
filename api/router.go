package api

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine) {
	router.POST("/storeEventData", StoreEventData)
	router.GET("/totalEventCountByType", TotalEventCountByType)
	router.GET("/eventCountByCameraId", EventCountByCameraID)
	router.GET("/totalEventCountByCameraId", TotalEventCountByCameraID) // TODO: check for validity of data
	router.GET("/getAllEventData", GetAllEventData)
}

// create interface - method to store data
// implement method for pulsar structure\
// interface for stote
// similarly for other handlers

// move to main.go
