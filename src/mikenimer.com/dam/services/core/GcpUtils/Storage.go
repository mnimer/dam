package Gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func DownloadFile(_bucket string, _name string) (*os.File, error) {
	//download file from GCS
	fileBytes, gcsErr := DownloadGcsFile(_bucket, _name)
	if gcsErr != nil {
		log.Fatal("Unable to download file `" + _bucket + ":" + _name + "`")
	}
	//save file to tmp dir
	filePath, err := WriteTempFile(_name, fileBytes)
	if err != nil {
		log.Fatal("Unable to save file `" + _bucket + ":" + _name + "`")
	}

	return filePath, err
}

func DownloadGcsFile(bucket, name string) ([]byte, error) {
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
