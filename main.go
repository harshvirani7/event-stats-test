package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/api"
	"github.com/harshvirani7/event-stats-test/storage"
)

func main() {
	// Initialize Redis client
	storage.InitRedisClient()

	// Initialize Gin router
	router := gin.Default()

	// Initialize routes
	api.InitRoutes(router)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// get data from pulsar and modify and store in redis

// query param - event type - return total event count
// GET - /countByEventType/{eventType}
// for each esn - count of total no of events - qp - event type
// GET - /totalEventsByESN/{eventType} // total events of that type for each esn
// count of event by esn id
// GET /totalEvents/{ESNID} -> int
// each event count for the esn
// GET /esnSummary/{ESNID}

// format data from pulsar and store to redis
// APIs will fetch data from redis store

// cameraId, timestamp, eventype
// POST API - workflow document - process and store relevant data in redis
