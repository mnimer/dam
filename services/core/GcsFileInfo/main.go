// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", RequestHandler)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// HelloPubSub receives and processes a Pub/Sub push message.
func RequestHandler(w http.ResponseWriter, r *http.Request) {

	//get NAME parameters
	names, ok := r.URL.Query()["name"]
	if !ok || len(names) != 1 {
		log.Fatal("Single Url Param 'name' is missing")
		return
	}
	//get Bucket parameters
	buckets, ok := r.URL.Query()["bucket"]
	if !ok || len(buckets) != 1 {
		log.Fatal("Single Url Param 'bucket' is missing")
		return
	}

	log.Printf("Core GcsFile Info | gs://%s/%s", buckets[0], names[0] )

	attr,err := getFileInfo(buckets[0], names[0])
	if err != nil { /* todo */ }

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attr)

}

func getFileInfo(bucket string, name string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}


	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(name)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Bucket: %v\n", attrs.Bucket)
	log.Printf("CacheControl: %v\n", attrs.CacheControl)
	log.Printf("ContentDisposition: %v\n", attrs.ContentDisposition)
	log.Printf("ContentEncoding: %v\n", attrs.ContentEncoding)
	log.Printf("ContentLanguage: %v\n", attrs.ContentLanguage)
	log.Printf("ContentType: %v\n", attrs.ContentType)
	log.Printf("Crc32c: %v\n", attrs.CRC32C)
	log.Printf("Generation: %v\n", attrs.Generation)
	log.Printf("KmsKeyName: %v\n", attrs.KMSKeyName)
	log.Printf("Md5Hash: %v\n", attrs.MD5)
	log.Printf("MediaLink: %v\n", attrs.MediaLink)
	log.Printf("Metageneration: %v\n", attrs.Metageneration)
	log.Printf("Name: %v\n", attrs.Name)
	log.Printf("Size: %v\n", attrs.Size)
	log.Printf("StorageClass: %v\n", attrs.StorageClass)
	log.Printf("TimeCreated: %v\n", attrs.Created)
	log.Printf("Updated: %v\n", attrs.Updated)
	log.Printf("Event-based hold enabled? %t\n", attrs.EventBasedHold)
	log.Printf("Temporary hold enabled? %t\n", attrs.TemporaryHold)
	log.Printf("Retention expiration time %v\n", attrs.RetentionExpirationTime)
	log.Print("\n\nMetadata\n")
	for key, value := range attrs.Metadata {
		log.Printf("\t%v = %v\n", key, value)
	}

	return attrs, nil
}
