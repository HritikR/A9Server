package main

import (
	"fmt"
	"log"
	"net/http"

	pppp "github.com/HritikR/A9Server/lib"
)

func main() {
	connection, err := pppp.InitiateConnection()
	if err != nil {
		log.Fatalf("Failed to initiate connection: %v", err)
	}
	defer connection.Close()

	stream := connection.RequestVideoStream()

	// Serve video stream over HTTP
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=--FRAMEBOUNDARY")
		for frame := range stream {
			// Write each frame with boundary and headers
			fmt.Fprintf(w, "--FRAMEBOUNDARY\r\nContent-Type: image/jpeg\r\n\r\n")
			w.Write(frame.Frame)
			fmt.Fprintf(w, "\r\n")
			w.(http.Flusher).Flush() // Ensure frame is sent immediately
		}
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
