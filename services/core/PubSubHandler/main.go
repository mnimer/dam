
// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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


// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type Data struct {
	Kind string 					`json:"kind"`
	Id string						`json:"id"`
	Selflink string					`json:"selflink"`
	Name string						`json:"name"`
	Bucket string					`json:"bucket"`
	Generation string				`json:"generation"`
	Metageneration string			`json:"metageneration"`
	ContentType string				`json:"contentType"`
	TimeCreated string				`json:"timeCreated"`
	Updated string					`json:"updated"`
	StorageClass string				`json:"storageClass"`
	TimeStorageClassUpdated  string	`json:"timeStorageClassUpdated"`
	Size string						`json:"size"`
	Md5Hash string					`json:"md5Hash"`
	MediaLink string				`json:"mediaLink"`
	Crc32c string					`json:"crc32c"`
	Etag string						`json:"etag"`
}


// MessageHandler receives and processes a Pub/Sub push message.
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	var m PubSubMessage
	var msgBody Data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//parse the outer messageg
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//parse the message body
	if err2 := json.Unmarshal(m.Message.Data, &msgBody); err2 != nil {
		log.Printf("json2.Unmarshal: %v", err2)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//fmt.Printf("[GCS] Kind: %s, Name: %s, Bucket: %s, ContentType: %s, Id: %s", msgBody.Kind, msgBody.Name, msgBody.Bucket, msgBody.ContentType, msgBody.Id)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgBody)
}