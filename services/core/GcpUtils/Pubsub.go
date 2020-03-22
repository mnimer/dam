package Gcp

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type PubSubData struct {
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


func ParsePubSubMessage(w http.ResponseWriter, r *http.Request) (PubSubData, error) {
	var m PubSubMessage
	var msgBody PubSubData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return PubSubData{}, err
	}else{
		println("body----")
		println(string(body))
	}
	//parse the outer messageg
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return PubSubData{}, err
	}
	//parse the message body
	if err2 := json.Unmarshal(m.Message.Data, &msgBody); err2 != nil {
		log.Printf("json2.Unmarshal: %v", err2)
		log.Println(m.Message.Data)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return PubSubData{}, err
	}
	return msgBody, nil
}



func SendToPubSub( topic string, msg []byte ) (error) {
	projectId, err := GetProjectId()
	if( err != nil ){
		return err
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topic)
	t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	return nil
}

