package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Define the HTTP handler function for your API endpoint
	http.HandleFunc("/api/basic-text", func(w http.ResponseWriter, r *http.Request) {
		// Set the content type header
		w.Header().Set("Content-Type", "text/plain")

		// Write the response body
		fmt.Fprintf(w, "This is some basic text returned by the API.")
	})

	// Start the HTTP server
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

// use gin
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
