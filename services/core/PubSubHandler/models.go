package main


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
