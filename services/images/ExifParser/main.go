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
	log.Printf("Image EXIF Parser | gs://%s/%s", _bucket, _name)

	// Process File
	filePath, err := gcp.DownloadFile(_bucket, _name)
	if err != nil {
		log.Fatal(err.Error())
	}

	//delete tmp file at end of method call
	defer gcp.DeleteTempFile(filePath.Name())

	// Parse
	metadata, err := ParseExif(filePath.Name())
	if err != nil {
		log.Fatal(err)
	}

	//break into multiple parts and send over PubSub
	FormatAndSaveProperties(_bucket, _name, metadata)

	//return if called via http
	w.WriteHeader(http.StatusOK)
	m := make(map[string]interface{})
	m["Name"] = _name
	m["Bucket"] = _bucket
	m["Metadata"] = metadata
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// Format the exif metadata into multiple smaller objects then
// Send the metadata to Pub/Sub
func FormatAndSaveProperties(bucket, name string, metadata map[string]interface{}) {

	primaryTags, err1 := FormatPrimary(bucket, name, metadata)
	if err1 == nil {
		primaryJson, err2 := json.Marshal(primaryTags)
		if err2 == nil {
			gcp.SaveMetadataFile(bucket, name, "exif.primary.json", primaryJson)
		} else {
			log.Println("Error parsing exif PrimaryTags | " + err2.Error())
		}
	}

	exifTags, err3 := FormatExifTags(bucket, name, metadata)
	if err3 == nil {
		exifJson, err4 := json.Marshal(exifTags)
		if err4 == nil {
			gcp.SaveMetadataFile(bucket, name, "exif.tags.json", exifJson)
		} else {
			log.Println("Error parsing exif Tags | " + err4.Error())
		}
	}

	gpsTags, err5 := FormatGps(bucket, name, metadata)
	if err5 == nil {
		gpsJson, err6 := json.Marshal(gpsTags)
		if err6 == nil {
			gcp.SaveMetadataFile(bucket, name, "exif.gps.json", gpsJson)
		} else {
			log.Println("Error parsing gps Tags | " + err6.Error())
		}
	}

	geoTags, err7 := FormatGeo(bucket, name, metadata)
	if err7 == nil {
		geoJson, err8 := json.Marshal(geoTags)
		if err8 == nil {
			gcp.SaveMetadataFile(bucket, name, "geo.json", geoJson)
		} else {
			log.Println("Error parsing geo Tags | " + err8.Error())
		}
	}
}

// Parse EXIF data out of the supported files, using a wrapper around the popular EXIFTOOL
// This method will return a KV list of all tags in the file.  We will break these into
// different files later, in the Format* methods.
func ParseExif(file string) (map[string]interface{}, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil, err
	}
	defer et.Close()
	fileInfos := et.ExtractMetadata(file)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}
	}

	delete(fileInfos[0].Fields, "Directory")
	delete(fileInfos[0].Fields, "FileName")
	delete(fileInfos[0].Fields, "SourceFile")
	return fileInfos[0].Fields, nil
}

// Pull out the key Primary tags that are common in
// If none of the known keys exists, will return nil so we do not save a file.
func FormatPrimary(bucket string, name string, metadata map[string]interface{}) (ExifPrimary, error) {

	primaryTags := ExifPrimary{bucket, name, time.Now(), PrimaryTags{}}

	keysFound := false
	primaryKeys := reflect.TypeOf(primaryTags)
	primaryMap := make(map[string]interface{})
	for i := 0; i < primaryKeys.NumField(); i++ {
		name := primaryKeys.Field(i).Name
		if metadata[name] != nil {
			keysFound = true
			primaryMap[name] = metadata[name]
		}
	}
	mapstructure.Decode(primaryMap, &primaryTags)

	if !keysFound {
		e := errors.New("No keys found")
		return primaryTags, e
	}
	return primaryTags, nil
}

// Pull out all of the exif tags (excluding primary and gps)
// If none of the known keys exists, will return nil so we do not save a file.
func FormatExifTags(bucket string, name string, metadata map[string]interface{}) (ExifTags, error) {
	//Initialize
	primaryTags := ExifPrimary{}
	gpsTags := ExifGps{}
	exifTags := ExifTags{bucket, name, time.Now(), make(map[string]interface{})}

	//Save the rest in a open map object
	keysFound := false
	primaryKeys := reflect.TypeOf(primaryTags)
	gpsKeys := reflect.TypeOf(gpsTags)
	for k, v := range metadata {
		pField, pBool := primaryKeys.FieldByName(k)
		gField, gBool := gpsKeys.FieldByName(k)

		if (pField.Name == k && !pBool) || (gField.Name == k && !gBool) {
			if k == "SourceFile" || k == "Directory" {
				//skip, it is pointing to the cloud run tmp dir, useless
			} else {
				keysFound = true
				exifTags.Metadata[k] = v
			}
		} else {
			//do nothing, because it is already in one of the other metadata structs
			//fmt.Printf("exists - key[%s] value[%s]\n", k, v)
		}
	}

	if !keysFound {
		e := errors.New("No keys found")
		return exifTags, e
	}
	return exifTags, nil
}

//Pull out the GPS specific tags from the exif file
// If none of the known keys exists, will return nil so we do not save a file.
func FormatGps(bucket string, name string, metadata map[string]interface{}) (ExifGps, error) {
	//Initialize
	gpsTags := ExifGps{bucket, name, time.Now(), GpsTags{}}

	keysFound := false
	gpsKeys := reflect.TypeOf(gpsTags)
	gpsMap := make(map[string]interface{})
	for i := 0; i < gpsKeys.NumField(); i++ {
		name := gpsKeys.Field(i).Name
		if metadata[name] != nil {
			keysFound = true
			gpsMap[name] = metadata[name]
		}
	}
	gpsStruct := &gpsTags
	mapstructure.Decode(gpsMap, gpsStruct)

	if !keysFound {
		e := errors.New("No keys found")
		return gpsTags, e
	}
	return gpsTags, nil
}

// Extract the  Lat/Lng properties from the GPS keys and save in a generic way
// We are going to save this in a generic way because there could be other ways to figure out the geo props of a file, not just exif
// If none of the known keys exists, will return nil so we do not save a file.
func FormatGeo(bucket string, name string, metadata map[string]interface{}) (ExifGeo, error) {
	//Initialize
	gpsTags := ExifGps{bucket, name, time.Now(), GpsTags{}}
	geoTags := ExifGeo{bucket, name, time.Now(), Geo{}}

	keysFound := false
	gpsKeys := reflect.TypeOf(gpsTags)
	gpsMap := make(map[string]interface{})
	for i := 0; i < gpsKeys.NumField(); i++ {
		name := gpsKeys.Field(i).Name
		if metadata[name] != nil {
			//pull out the GPS Lat/Lng  columns and save as generic lat/lng properties
			if name == "GPSLongitude" || name == "GPSLatitude" {
				keysFound = true
				val := parseCoordString(metadata[name].(string))
				gpsMap[name] = val
				if name == "GPSLongitude" {
					geoTags.Metadata.Longitude = val
				} else if name == "GPSLatitude" {
					geoTags.Metadata.Latitude = val
				}
			}
		}
	}
	gpsStruct := &gpsTags
	mapstructure.Decode(gpsMap, gpsStruct)

	if !keysFound {
		e := errors.New("No keys found")
		return geoTags, e
	}
	return geoTags, nil
}

// Parse the lat/lng strings returned from the ExifTool and parse into a valid float number
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
