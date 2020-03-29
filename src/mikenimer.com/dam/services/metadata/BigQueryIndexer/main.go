// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/barasher/go-exiftool"
	"github.com/mitchellh/mapstructure"
	"log"
	gcp "mikenimer.com/services/core/GcpUtils"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", RequestHandler)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// MessageHandler receives and processes a Pub/Sub push message.
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	///////
	//Validate Arguments
	msgBody, pErr := gcp.ParsePubSubMessage(w, r)
	if pErr != nil {
		log.Fatal(pErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(pErr.Error()))
		return
	}

	_bucket := msgBody.Bucket
	_name := msgBody.Name
	log.Printf("BigQuery Metadata Parser | gs://%s/%s", _bucket, _name)
}
