package main

import (
	"cloud.google.com/go/vision/apiv1"
	"context"
	"encoding/json"
	"fmt"
	vision3 "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"log"
	gcp "mikenimer.com/services/core/GcpUtils"
	"net/http"
	"os"
	"time"
)

type LabelsApiResult struct{
	Bucket	string			`json:"bucket"`
	Name 	string			`json:"name"`
	DateTime time.Time		`json:"datetime"`
	Metadata []*vision3.EntityAnnotation	`json:"metadata"`
}

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
	if( pErr != nil ){
		log.Fatal(pErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(pErr.Error()))
		return
	}

	_bucket:=msgBody.Bucket
	_name:=msgBody.Name
	log.Printf("Image VisionAPI Parser | gs://%s/%s", _bucket, _name )


	apiResults, err := callVisionApi(_bucket, _name)

	if( err != nil ){
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	} else {
		if (len(apiResults) > 0) {
			resultsJson, err2 := json.Marshal(apiResults)
			if (err2 == nil) {
				gcp.SaveMetadataFile(_bucket, _name, "visionapi.labels.json", resultsJson)
			} else {
				log.Println("Error parsing VisionAPI Labels | " + err2.Error())
			}
		}
	}

	//return if called via http
	w.WriteHeader(http.StatusOK)
	m :=  LabelsApiResult{}
	m.Name = _name
	m.Bucket = _bucket
	m.DateTime = time.Now()
	m.Metadata = apiResults
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}



func callVisionApi(bucket string, name string) ([]*vision3.EntityAnnotation, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		//log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}
	defer client.Close()

	uri := fmt.Sprintf("gs://%s/%s", bucket, name)
	image := vision.NewImageFromURI(uri)
	if image == nil {
		//log.Fatalf("Failed to create image: %v", err)
		return nil, err
	}

	results, err := client.DetectLabels(ctx, image, nil, 20)
	if err != nil {
		//log.Fatalf("Failed to detect labels: %v", err)
		return nil, err
	}
	return results, nil
}