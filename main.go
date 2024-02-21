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
