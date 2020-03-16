// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/barasher/go-exiftool"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
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


	log.Printf("Image EXIF Parser | gs://%s/%s", buckets[0], names[0] )

	//get local var
	file := files[0]
	name := names[0]
	bucket := buckets[0]

	///////
	// Process File

	//download file from GCS
	fileBytes, gcsErr := DownloadFile(bucket, name, file)
	if gcsErr != nil {
		log.Fatal("Unable to download file `" + file + "`")
	}
	//save file to tmp dir
	filePath, err := WriteTempFile(name, fileBytes)
	if err != nil {
		log.Fatal("Unable to save file `" + file + "`")
	}
	//delete tmp file at end of method call
	defer DeleteTempFile(filePath.Name())

	// Query()["files"] will return an array of items, we only want the single item.
	metadata, err := ParseExif(filePath.Name())
	if err != nil {
		log.Fatal(err)
	}

	exif := FormatResults(bucket, name, file, metadata)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exif)
}

func DownloadFile(bucket, name, file string) ([]byte, error) {
	// [START download_file]
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	//filePart := strings.TrimLeft(file, bucket +"/")
	//object := bucket + ":" + name
	ctx, cancel := context.WithTimeout(ctx, time.Second*300) //todo: make this timeout configurable
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
	nameParts := strings.SplitAfter(name, ".")
	ext := nameParts[len(nameParts)-1]

	tmpFile, err := ioutil.TempFile(os.TempDir(), "exif-*."+ext)
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	//fmt.Println("Created Tmp File: " + tmpFile.Name())

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

func ParseExif(file string) (map[string]interface{}, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil, errors.New("Error initializing EXIFTOOL")
	}
	defer et.Close()
	fileInfos := et.ExtractMetadata(file)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}
	}

	return fileInfos[0].Fields, nil
}

func FormatResults(bucket string, name string, file string, metadata map[string]interface{}) interface{} {
	//Initialize
	e := Exif{file, bucket, name, ExifMetadata{PrimaryTags{}, GpsTags{}, make(map[string]interface{})}}

	primaryKeys := reflect.TypeOf(e.Metadata.PrimaryTags)
	primaryMap := make(map[string]interface{})
	for i := 0; i < primaryKeys.NumField(); i++ {
		name := primaryKeys.Field(i).Name
		if metadata[name] != nil {
			primaryMap[name] = metadata[name]
		}
	}
	mapstructure.Decode(primaryMap, &e.Metadata.PrimaryTags)

	gpsKeys := reflect.TypeOf(e.Metadata.GpsTags)
	gpsMap := make(map[string]interface{})
	for i := 0; i < gpsKeys.NumField(); i++ {
		name := gpsKeys.Field(i).Name
		if metadata[name] != nil {
			if name == "GPSLongitude" || name == "GPSLatitude" {
				gpsMap[name] = fmt.Sprintf("%f", parseCoordString(metadata[name].(string)))
			} else {
				gpsMap[name] = metadata[name]
			}
		}
	}

	gpsStruct := &e.Metadata.GpsTags
	mapstructure.Decode(gpsMap, gpsStruct)
	//e.Metadata.GpsTags = gpsStruct

	//Save the rest in a open map object
	for k, v := range metadata {
		if primaryMap[k] == nil && gpsMap[k] == nil {
			if( k == "SourceFile" || k == "Directory"){
				//skip, it is pointing to the cloud run tmp dir, useless
			}else {
				e.Metadata.ExifTags[k] = v;
			}
		} else {
			//do nothing, because it is already in one of the other metadata structs
			//fmt.Printf("exists - key[%s] value[%s]\n", k, v)
		}
	}

	return e
}

func parseCoordString(gpsLatLng string) float64 {
	clean := strings.ReplaceAll(gpsLatLng, "'", "")
	clean = strings.ReplaceAll(clean, "\"", "")
	chunks := strings.Split(clean, " ")
	hours, _ := strconv.ParseFloat(strings.TrimSpace(chunks[0]), 64)
	minutes, _ := strconv.ParseFloat(strings.TrimSpace(chunks[2]), 64)
	seconds, _ := strconv.ParseFloat(strings.TrimSpace(chunks[3]), 64)
	latlng := hours + (minutes / 60) + (seconds / 3600)

	direction := strings.Split(gpsLatLng, " ")
	if direction[len(direction)-1] == "S" { // N is "+", S is "-"
		return latlng * -1
	} else if direction[len(direction)-1] == "W" { // E is "+", W is "-"
		return latlng * -1
	} else {
		return latlng
	}

}
