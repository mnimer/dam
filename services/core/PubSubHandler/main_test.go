package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func TestObjectFinalizeMsg(t *testing.T) {
	// Create a request to pass to our handler.
	var testMsg string = "{\"message\":{\"attributes\":{\"bucketId\":\"mikenimer-dam-playground-content\",\"eventTime\":\"2020-03-10T00:37:25.931866Z\",\"eventType\":\"OBJECT_FINALIZE\",\"notificationConfig\":\"projects/_/buckets/mikenimer-dam-playground-content/notificationConfigs/1\",\"objectGeneration\":\"1583800645932016\",\"objectId\":\"AlaskanGlacier.jpg\",\"payloadFormat\":\"JSON_API_V1\"},\"data\":\"ewogICJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwKICAiaWQiOiAibWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvQWxhc2thbkdsYWNpZXIuanBnLzE1ODM4MDA2NDU5MzIwMTYiLAogICJzZWxmTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby9BbGFza2FuR2xhY2llci5qcGciLAogICJuYW1lIjogIkFsYXNrYW5HbGFjaWVyLmpwZyIsCiAgImJ1Y2tldCI6ICJtaWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudCIsCiAgImdlbmVyYXRpb24iOiAiMTU4MzgwMDY0NTkzMjAxNiIsCiAgIm1ldGFnZW5lcmF0aW9uIjogIjEiLAogICJjb250ZW50VHlwZSI6ICJpbWFnZS9qcGVnIiwKICAidGltZUNyZWF0ZWQiOiAiMjAyMC0wMy0xMFQwMDozNzoyNS45MzFaIiwKICAidXBkYXRlZCI6ICIyMDIwLTAzLTEwVDAwOjM3OjI1LjkzMVoiLAogICJzdG9yYWdlQ2xhc3MiOiAiU1RBTkRBUkQiLAogICJ0aW1lU3RvcmFnZUNsYXNzVXBkYXRlZCI6ICIyMDIwLTAzLTEwVDAwOjM3OjI1LjkzMVoiLAogICJzaXplIjogIjcwOTYzMzIiLAogICJtZDVIYXNoIjogIlljUWM2YjNSaVozVVhGcEducFRRV2c9PSIsCiAgIm1lZGlhTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9kb3dubG9hZC9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby9BbGFza2FuR2xhY2llci5qcGc/Z2VuZXJhdGlvbj0xNTgzODAwNjQ1OTMyMDE2JmFsdD1tZWRpYSIsCiAgImNyYzMyYyI6ICJpZktWMFE9PSIsCiAgImV0YWciOiAiQ1BDbmxaZlZqdWdDRUFFPSIKfQo=\",\"messageId\":\"1037428102269218\",\"message_id\":\"1037428102269218\",\"publishTime\":\"2020-03-10T00:37:26.069Z\",\"publish_time\":\"2020-03-10T00:37:26.069Z\"},\"subscription\":\"projects/mikenimer-dam-playground/subscriptions/dam-pubsubhandler\"}"

	reader := strings.NewReader(testMsg)
	req, err := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}


	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MessageHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}


	var msgBody Data
	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &msgBody); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		return
	}

	var msgExpected Data
	// Check the response body is what we expect.
	const expected string = `{"kind":"storage#object","id":"mikenimer-dam-playground-content/AlaskanGlacier.jpg/1583800645932016","selflink":"https://www.googleapis.com/storage/v1/b/mikenimer-dam-playground-content/o/AlaskanGlacier.jpg","name":"AlaskanGlacier.jpg","bucket":"mikenimer-dam-playground-content","generation":"1583800645932016","metageneration":"1","contentType":"image/jpeg","timeCreated":"2020-03-10T00:37:25.931Z","updated":"2020-03-10T00:37:25.931Z","storageClass":"STANDARD","timeStorageClassUpdated":"2020-03-10T00:37:25.931Z","size":"7096332","md5Hash":"YcQc6b3RiZ3UXFpGnpTQWg==","mediaLink":"https://www.googleapis.com/download/storage/v1/b/mikenimer-dam-playground-content/o/AlaskanGlacier.jpg?generation=1583800645932016\u0026alt=media","crc32c":"ifKV0Q==","etag":"CPCnlZfVjugCEAE="}`
	//parse the expected message
	if err2 := json.Unmarshal( []byte(expected), &msgExpected); err2 != nil {
		log.Printf("json.Unmarshal: %v", err2)
		return
	}

	assert := assert.New(t)
	assert.Equal(rr.Code, 200, "Incorrect Header")
	assert.Equal(rr.Header().Get("Content-Type"), "application/json", "Incorrect Header")
	assert.Equal(msgBody.Id, msgExpected.Id, "ID does not match")
	assert.Equal(msgBody.Kind, msgExpected.Kind, "Kind does not match")
	assert.Equal(msgBody.Name, msgExpected.Name, "Name does not match")
	assert.Equal(msgBody.Bucket, msgExpected.Bucket, "Bucket does not match")
	assert.Equal(msgBody.ContentType, msgExpected.ContentType, "ContentType does not match")

}

func TestFileInSubFolderMsg(t *testing.T) {
	// Create a request to pass to our handler.
	var testMsg string = `{"message":{"attributes":{"bucketId":"mikenimer-dam-playground-content","eventTime":"2020-03-10T04:48:23.174091Z","eventType":"OBJECT_FINALIZE","notificationConfig":"projects/_/buckets/mikenimer-dam-playground-content/notificationConfigs/1","objectGeneration":"1583815703174289","objectId":"folder1/AlaskanGlacier.jpg","payloadFormat":"JSON_API_V1"},"data":"ewogICJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwKICAiaWQiOiAibWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvZm9sZGVyMS9BbGFza2FuR2xhY2llci5qcGcvMTU4MzgxNTcwMzE3NDI4OSIsCiAgInNlbGZMaW5rIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3N0b3JhZ2UvdjEvYi9taWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudC9vL2ZvbGRlcjElMkZBbGFza2FuR2xhY2llci5qcGciLAogICJuYW1lIjogImZvbGRlcjEvQWxhc2thbkdsYWNpZXIuanBnIiwKICAiYnVja2V0IjogIm1pa2VuaW1lci1kYW0tcGxheWdyb3VuZC1jb250ZW50IiwKICAiZ2VuZXJhdGlvbiI6ICIxNTgzODE1NzAzMTc0Mjg5IiwKICAibWV0YWdlbmVyYXRpb24iOiAiMSIsCiAgImNvbnRlbnRUeXBlIjogImltYWdlL2pwZWciLAogICJ0aW1lQ3JlYXRlZCI6ICIyMDIwLTAzLTEwVDA0OjQ4OjIzLjE3NFoiLAogICJ1cGRhdGVkIjogIjIwMjAtMDMtMTBUMDQ6NDg6MjMuMTc0WiIsCiAgInN0b3JhZ2VDbGFzcyI6ICJTVEFOREFSRCIsCiAgInRpbWVTdG9yYWdlQ2xhc3NVcGRhdGVkIjogIjIwMjAtMDMtMTBUMDQ6NDg6MjMuMTc0WiIsCiAgInNpemUiOiAiNzA5NjMzMiIsCiAgIm1kNUhhc2giOiAiWWNRYzZiM1JpWjNVWEZwR25wVFFXZz09IiwKICAibWVkaWFMaW5rIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL2Rvd25sb2FkL3N0b3JhZ2UvdjEvYi9taWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudC9vL2ZvbGRlcjElMkZBbGFza2FuR2xhY2llci5qcGc/Z2VuZXJhdGlvbj0xNTgzODE1NzAzMTc0Mjg5JmFsdD1tZWRpYSIsCiAgImNyYzMyYyI6ICJpZktWMFE9PSIsCiAgImV0YWciOiAiQ0pHNWdxT05qK2dDRUFFPSIKfQo=","messageId":"1037613170822282","message_id":"1037613170822282","publishTime":"2020-03-10T04:48:23.385Z","publish_time":"2020-03-10T04:48:23.385Z"},"subscription":"projects/mikenimer-dam-playground/subscriptions/dam-pubsubhandler"}`

	reader := strings.NewReader(testMsg)
	req, err := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}


	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MessageHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}


	var msgBody Data
	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &msgBody); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		return
	}

	var msgExpected Data
	// Check the response body is what we expect.
	const expected string = `{"message":{"attributes":{"bucketId":"mikenimer-dam-playground-content","eventTime":"2020-03-10T04:48:23.174091Z","eventType":"OBJECT_FINALIZE","notificationConfig":"projects/_/buckets/mikenimer-dam-playground-content/notificationConfigs/1","objectGeneration":"1583815703174289","objectId":"folder1/AlaskanGlacier.jpg","payloadFormat":"JSON_API_V1"},"data":"ewogICJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwKICAiaWQiOiAibWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvZm9sZGVyMS9BbGFza2FuR2xhY2llci5qcGcvMTU4MzgxNTcwMzE3NDI4OSIsCiAgInNlbGZMaW5rIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3N0b3JhZ2UvdjEvYi9taWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudC9vL2ZvbGRlcjElMkZBbGFza2FuR2xhY2llci5qcGciLAogICJuYW1lIjogImZvbGRlcjEvQWxhc2thbkdsYWNpZXIuanBnIiwKICAiYnVja2V0IjogIm1pa2VuaW1lci1kYW0tcGxheWdyb3VuZC1jb250ZW50IiwKICAiZ2VuZXJhdGlvbiI6ICIxNTgzODE1NzAzMTc0Mjg5IiwKICAibWV0YWdlbmVyYXRpb24iOiAiMSIsCiAgImNvbnRlbnRUeXBlIjogImltYWdlL2pwZWciLAogICJ0aW1lQ3JlYXRlZCI6ICIyMDIwLTAzLTEwVDA0OjQ4OjIzLjE3NFoiLAogICJ1cGRhdGVkIjogIjIwMjAtMDMtMTBUMDQ6NDg6MjMuMTc0WiIsCiAgInN0b3JhZ2VDbGFzcyI6ICJTVEFOREFSRCIsCiAgInRpbWVTdG9yYWdlQ2xhc3NVcGRhdGVkIjogIjIwMjAtMDMtMTBUMDQ6NDg6MjMuMTc0WiIsCiAgInNpemUiOiAiNzA5NjMzMiIsCiAgIm1kNUhhc2giOiAiWWNRYzZiM1JpWjNVWEZwR25wVFFXZz09IiwKICAibWVkaWFMaW5rIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL2Rvd25sb2FkL3N0b3JhZ2UvdjEvYi9taWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudC9vL2ZvbGRlcjElMkZBbGFza2FuR2xhY2llci5qcGc/Z2VuZXJhdGlvbj0xNTgzODE1NzAzMTc0Mjg5JmFsdD1tZWRpYSIsCiAgImNyYzMyYyI6ICJpZktWMFE9PSIsCiAgImV0YWciOiAiQ0pHNWdxT05qK2dDRUFFPSIKfQo=","messageId":"1037613170822282","message_id":"1037613170822282","publishTime":"2020-03-10T04:48:23.385Z","publish_time":"2020-03-10T04:48:23.385Z"},"subscription":"projects/mikenimer-dam-playground/subscriptions/dam-pubsubhandler"}`
	//parse the expected message
	if err2 := json.Unmarshal( []byte(expected), &msgExpected); err2 != nil {
		log.Printf("json.Unmarshal: %v", err2)
		return
	}

	assert := assert.New(t)
	assert.Equal(rr.Code, 200, "Incorrect Header")
	assert.Equal(rr.Header().Get("Content-Type"), "application/json", "Incorrect Header")
	assert.Equal(msgBody.Id, msgExpected.Id, "ID does not match")
	assert.Equal(msgBody.Kind, msgExpected.Kind, "Kind does not match")
	assert.Equal(msgBody.Name, msgExpected.Name, "Name does not match")
	assert.Equal(msgBody.Bucket, msgExpected.Bucket, "Bucket does not match")
	assert.Equal(msgBody.ContentType, msgExpected.ContentType, "ContentType does not match")

}




func TestCreateFolderMsg(t *testing.T) {
	// Create a request to pass to our handler.
	var testMsg string = `{"message":{"attributes":{"bucketId":"mikenimer-dam-playground-content","eventTime":"2020-03-10T04:40:05.625245Z","eventType":"OBJECT_FINALIZE","notificationConfig":"projects/_/buckets/mikenimer-dam-playground-content/notificationConfigs/1","objectGeneration":"1583815205625424","objectId":"folder1/","payloadFormat":"JSON_API_V1"},"data":"ewogICJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwKICAiaWQiOiAibWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvZm9sZGVyMS8vMTU4MzgxNTIwNTYyNTQyNCIsCiAgInNlbGZMaW5rIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3N0b3JhZ2UvdjEvYi9taWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudC9vL2ZvbGRlcjElMkYiLAogICJuYW1lIjogImZvbGRlcjEvIiwKICAiYnVja2V0IjogIm1pa2VuaW1lci1kYW0tcGxheWdyb3VuZC1jb250ZW50IiwKICAiZ2VuZXJhdGlvbiI6ICIxNTgzODE1MjA1NjI1NDI0IiwKICAibWV0YWdlbmVyYXRpb24iOiAiMSIsCiAgImNvbnRlbnRUeXBlIjogInRleHQvcGxhaW4iLAogICJ0aW1lQ3JlYXRlZCI6ICIyMDIwLTAzLTEwVDA0OjQwOjA1LjYyNVoiLAogICJ1cGRhdGVkIjogIjIwMjAtMDMtMTBUMDQ6NDA6MDUuNjI1WiIsCiAgInN0b3JhZ2VDbGFzcyI6ICJTVEFOREFSRCIsCiAgInRpbWVTdG9yYWdlQ2xhc3NVcGRhdGVkIjogIjIwMjAtMDMtMTBUMDQ6NDA6MDUuNjI1WiIsCiAgInNpemUiOiAiMTEiLAogICJtZDVIYXNoIjogImFwbkZkYXVIK01mUjdSNVM1K05Kemc9PSIsCiAgIm1lZGlhTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9kb3dubG9hZC9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby9mb2xkZXIxJTJGP2dlbmVyYXRpb249MTU4MzgxNTIwNTYyNTQyNCZhbHQ9bWVkaWEiLAogICJjcmMzMmMiOiAiWGtJK0R3PT0iLAogICJldGFnIjogIkNOQzg0cldMaitnQ0VBRT0iLAogICJldmVudEJhc2VkSG9sZCI6IGZhbHNlCn0K","messageId":"1037604744720480","message_id":"1037604744720480","publishTime":"2020-03-10T04:40:06.004Z","publish_time":"2020-03-10T04:40:06.004Z"},"subscription":"projects/mikenimer-dam-playground/subscriptions/dam-pubsubhandler"}`

	reader := strings.NewReader(testMsg)
	req, err := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}


	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MessageHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}


	var msgBody Data
	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &msgBody); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		return
	}

	var msgExpected Data
	// Check the response body is what we expect.
	const expected string = `{"kind":"storage#object","id":"mikenimer-dam-playground-content/AlaskanGlacier.jpg/1583800645932016","selflink":"https://www.googleapis.com/storage/v1/b/mikenimer-dam-playground-content/o/AlaskanGlacier.jpg","name":"AlaskanGlacier.jpg","bucket":"mikenimer-dam-playground-content","generation":"1583800645932016","metageneration":"1","contentType":"image/jpeg","timeCreated":"2020-03-10T00:37:25.931Z","updated":"2020-03-10T00:37:25.931Z","storageClass":"STANDARD","timeStorageClassUpdated":"2020-03-10T00:37:25.931Z","size":"7096332","md5Hash":"YcQc6b3RiZ3UXFpGnpTQWg==","mediaLink":"https://www.googleapis.com/download/storage/v1/b/mikenimer-dam-playground-content/o/AlaskanGlacier.jpg?generation=1583800645932016\u0026alt=media","crc32c":"ifKV0Q==","etag":"CPCnlZfVjugCEAE="}`
	//parse the expected message
	if err2 := json.Unmarshal( []byte(expected), &msgExpected); err2 != nil {
		log.Printf("json.Unmarshal: %v", err2)
		return
	}

	assert := assert.New(t)
	assert.Equal(rr.Code, 200, "Incorrect Header")
	assert.Equal(rr.Header().Get("Content-Type"), "application/json", "Incorrect Header")
	assert.Equal(msgBody.Id, msgExpected.Id, "ID does not match")
	assert.Equal(msgBody.Kind, msgExpected.Kind, "Kind does not match")
	assert.Equal(msgBody.Name, msgExpected.Name, "Name does not match")
	assert.Equal(msgBody.Bucket, msgExpected.Bucket, "Bucket does not match")
	assert.Equal(msgBody.ContentType, msgExpected.ContentType, "ContentType does not match")

}