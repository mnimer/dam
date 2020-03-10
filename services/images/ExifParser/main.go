// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/barasher/go-exiftool"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", MessageHandler)
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
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Image EXIF Parser")

	//get URL parameters
	files, ok := r.URL.Query()["file"]
	if !ok || len(files) != 1 {
		log.Println("Single Url Param 'file' is missing")
		return
	}
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

	//get local var
	file := files[0]
	name := names[0]
	bucket := buckets[0]

	//download file from GCS
	fileBytes, gcsErr := DownloadFile(bucket, name, file)
	if( gcsErr != nil ){
		log.Fatal("Unable to download file `" +file +"`")
	}
	//save file to tmp dir
	filePath, err := WriteTempFile(name, fileBytes)
	if( err != nil ){
		log.Fatal("Unable to save file `" +file +"`")
	}

	//delete tmp file at end of method call
	defer DeleteTempFile(filePath.Name())


	// Query()["files"] will return an array of items, we only want the single item.
	metadata, err := ParseExif(filePath.Name())
	if err != nil {
		log.Fatal(err)
	}
	println("metadata %s", metadata)

}



func DownloadFile(bucket, name, file string) ([]byte, error) {
	// [START download_file]
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	//filePart := strings.TrimLeft(file, bucket +"/")
	//object := bucket +":" +name
	ctx, cancel := context.WithTimeout(ctx, time.Second*60) //todo: make this timeout configurable
	defer cancel()
	rc, err := client.Bucket(bucket).Object(name).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}


//Save file Bytes to tmp dir
func WriteTempFile(name string, file []byte) (*os.File, error) {

	// Create our Temp File:  This will create a filename like /tmp/prefix-123456
	// We can use a pattern of "pre-*.txt" to get an extension like: /tmp/pre-123456.txt
	nameParts  := strings.SplitAfter(name, ".")
	ext := nameParts[len(nameParts)-1];

	tmpFile, err := ioutil.TempFile(os.TempDir(), "exif-*." +ext)
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	fmt.Println("Created Tmp File: " + tmpFile.Name())

	// Example writing to the file
	data := []byte(file)
	if _, err = tmpFile.Write(data); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpFile, nil
}

//Delete tmp file
func DeleteTempFile(filePath string) {
	// Remember to clean up the file afterwards
	os.Remove(filePath)
}


func ParseExif(file string) (exiftool.FileMetadata, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return "{}", errors.New("Error initializing EXIFTOOL")
	}
	defer et.Close()
	fileInfos := et.ExtractMetadata(file)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			fmt.Printf("[%v] %v\n", k, v)
		}
	}

	return fileInfos, nil
}
