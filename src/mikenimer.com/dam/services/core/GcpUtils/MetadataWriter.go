package Gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"log"
)

func SaveMetadataFile(bucket, name, fileName string, content []byte) int {
	//get hidden dir to store metadata
	dirName := GetMetadataFolder(name)

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
