package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a route for the API endpoint
	router.GET("/api/basic-text", func(c *gin.Context) {
		// Set the response status code and content type
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/plain")

		// Write the response body
		c.String(http.StatusOK, "This is some basic text returned by the API.")
	})

	// Run the server
	router.Run(":8080")
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
