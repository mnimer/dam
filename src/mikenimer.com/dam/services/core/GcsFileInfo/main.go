// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"cloud.google.com/go/storage"
	gcp "github.com/mnimer/dam/services/core/GcpUtils"
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

	msgBody, pErr := gcp.ParsePubSubMessage(w, r)
	if pErr != nil {
		println(pErr)
		w.WriteHeader(http.StatusInternalServerError)
		//w.Write(pErr)
	}

	_bucket := msgBody.Bucket
	_name := msgBody.Name
	log.Printf("Core GcsFile Info | gs://%s/%s", _bucket, _name)

	attr, error := getFileInfo(_bucket, _name)
	if error != nil {
		println("Error getting FileInfo")
		println(error)
		w.WriteHeader(http.StatusInternalServerError)
	}

	//Send to one of two paths, Content Parsing or Indexing Metadata
	ForwardToPubSubTopics(attr)

	w.WriteHeader(http.StatusOK)
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(attr)
}

func getFileInfo(bucket string, name string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	o := client.Bucket(bucket).Object(name)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	/**
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
	}**/

	return attrs, nil
}

func ForwardToPubSubTopics(attrs *storage.ObjectAttrs) {
	msg, err := json.Marshal(attrs)
	if err == nil && attrs != nil {
		//Split the topics into two branchs, once for .metadata hidden files, that need to be indexed
		if strings.HasPrefix(attrs.Name, ".metadata") {
			println("Send to gcs-metadata-handlers")
			println(string(msg))
			gcp.SendToPubSub("gcs-metadata-handlers", msg)
		} else {
			// Send to image  format specific topics
			if strings.HasPrefix(attrs.ContentType, "image/") {
				println("Send to gcs-image-handlers")
				println(string(msg))
				gcp.SendToPubSub("gcs-image-handlers", msg)
			} else if strings.HasPrefix(attrs.ContentType, "video/") {
				println("Send to gcs-video-handlers")
				println(string(msg))
				gcp.SendToPubSub("gcs-video-handlers", msg)
			}
		}
	}
}
