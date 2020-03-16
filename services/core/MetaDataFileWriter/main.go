// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

	//extract raw body
	_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//convert body to metadata object
	var m Metadata
	if err := json.Unmarshal(_body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//validate metadata
	if m.Name == "" {
		log.Fatal("Missing file name")
	}
	if m.Bucket == "" {
		log.Fatal("Missing file bucket")
	}
	if m.FileId == "" {
		log.Fatal("Missing file Id")
	}

	log.Printf("Core MetaData File Writer | gs://%s/%s", m.Bucket, m.Name )

	//Make sure we have metadata to process
	if len(m.Metadata) > 0 {
		//Set HIDDEN folder path
		var metadataDirName = ".metadata/" + m.Name

		// Extract Geo Data from json
		if m.Metadata["GpsTags"] != nil {
			geoTags, err := FormatFloatTypes(m.Metadata["GpsTags"].(map[string]interface{}), "GPSLatitude", "GPSLongitude")
			if err == nil {
				////////////////////////////////////////
				//save exif gps tag
				geoTagsJson, err := json.Marshal(geoTags)
				if err == nil {
					SaveMetadataFile(m.Bucket, metadataDirName, "exif.geo.json", geoTagsJson)
				}

				////////////////////////////////////////
				//save 2nd file with lat/lng only data
				if geoTags["GPSLatitude"] != nil && geoTags["GPSLongitude"] != nil {
					geo := Geo{geoTags["GPSLatitude"].(float64), geoTags["GPSLongitude"].(float64)}
					geoJson, err := json.Marshal(geo)
					if err == nil {
						SaveMetadataFile(m.Bucket, metadataDirName, "geo.json", geoJson)
					}
				}
			} else {
				log.Println("Error converting GEO keys")
			}
		}

		////////////////////////////////////////
		// Extract jpeg tags
		if m.Metadata["JpegTags"] != nil {
			jpgTagsJson, err := json.Marshal(m.Metadata["JpegTags"])
			if err == nil { /*todo*/
				SaveMetadataFile(m.Bucket, metadataDirName, "exif.jpg.json", jpgTagsJson)
			}
		}
	}

}

//is key in array
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

/**
func GetProjectId() string {
	url := "http://metadata.google.internal/computeMetadata/v1/project/project-id";
	resp, err := http.Get(url)
	if err != nil {
		// todo handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// todo handle error
	}
	return string(body)
} **/

func SaveMetadataFile(bucket, dirName, fileName string, content []byte) int {
	//setup client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println("Error Creating GCS client")
	}

	//Save metadata in 'hidden' <bucket>/.metadata/<file name> GCS folder
	object := dirName + "/" + fileName
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

	n, err := wc.Write(content)
	defer wc.Close()

	if err != nil {
		log.Println("Error Saving Metadata File")
	}
	return n
}

func FormatFloatTypes(tags map[string]interface{}, floatKeys ...string) (map[string]interface{}, error) {
	for i := 0; i < len(floatKeys); i++ {
		fKey := floatKeys[i]
		if tags[fKey] != nil {
			fVal, err := strconv.ParseFloat(tags[fKey].(string), 64)
			if err == nil {
				tags[fKey] = fVal
			}
		}
	}
	return tags, nil
}
