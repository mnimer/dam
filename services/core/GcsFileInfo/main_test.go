package main

import (
	"cloud.google.com/go/storage"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//var fijiTestFile string = "{\n  \"insertId\": \"5e7055510008bb0357a0c012\",\n  \"jsonPayload\": {\n    \"message\": {\n      \"data\": \"ewogICJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwKICAiaWQiOiAibWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvdGVzdC9leGlmX2pwZ19GdWppZmlsbV9GaW5lUGl4NjkwMFpPT00uanBnLzE1ODQ0MTk5MjA5MTA2NTYiLAogICJzZWxmTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby90ZXN0JTJGZXhpZl9qcGdfRnVqaWZpbG1fRmluZVBpeDY5MDBaT09NLmpwZyIsCiAgIm5hbWUiOiAidGVzdC9leGlmX2pwZ19GdWppZmlsbV9GaW5lUGl4NjkwMFpPT00uanBnIiwKICAiYnVja2V0IjogIm1pa2VuaW1lci1kYW0tcGxheWdyb3VuZC1jb250ZW50IiwKICAiZ2VuZXJhdGlvbiI6ICIxNTg0NDE5OTIwOTEwNjU2IiwKICAibWV0YWdlbmVyYXRpb24iOiAiMSIsCiAgImNvbnRlbnRUeXBlIjogImltYWdlL2pwZWciLAogICJ0aW1lQ3JlYXRlZCI6ICIyMDIwLTAzLTE3VDA0OjM4OjQwLjkxMFoiLAogICJ1cGRhdGVkIjogIjIwMjAtMDMtMTdUMDQ6Mzg6NDAuOTEwWiIsCiAgInN0b3JhZ2VDbGFzcyI6ICJTVEFOREFSRCIsCiAgInRpbWVTdG9yYWdlQ2xhc3NVcGRhdGVkIjogIjIwMjAtMDMtMTdUMDQ6Mzg6NDAuOTEwWiIsCiAgInNpemUiOiAiNDI3OCIsCiAgIm1kNUhhc2giOiAiYTlKTjMxQ1V6Q3lNNm80UnM1ZlV6QT09IiwKICAibWVkaWFMaW5rIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL2Rvd25sb2FkL3N0b3JhZ2UvdjEvYi9taWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudC9vL3Rlc3QlMkZleGlmX2pwZ19GdWppZmlsbV9GaW5lUGl4NjkwMFpPT00uanBnP2dlbmVyYXRpb249MTU4NDQxOTkyMDkxMDY1NiZhbHQ9bWVkaWEiLAogICJjcmMzMmMiOiAicFpxa09RPT0iLAogICJldGFnIjogIkNNQ3l2cFRZb09nQ0VBRT0iCn0K\",\n      \"message_id\": \"1049793051698819\",\n      \"messageId\": \"1049793051698819\",\n      \"publishTime\": \"2020-03-17T04:42:56.341Z\",\n      \"publish_time\": \"2020-03-17T04:42:56.341Z\",\n      \"attributes\": {\n        \"objectId\": \"test/exif_jpg_Fujifilm_FinePix6900ZOOM.jpg\",\n        \"objectGeneration\": \"1584419920910656\",\n        \"eventTime\": \"2020-03-17T04:42:55.990279Z\",\n        \"eventType\": \"OBJECT_DELETE\",\n        \"payloadFormat\": \"JSON_API_V1\",\n        \"notificationConfig\": \"projects/_/buckets/mikenimer-dam-playground-content/notificationConfigs/1\",\n        \"overwrittenByGeneration\": \"1584420175990560\",\n        \"bucketId\": \"mikenimer-dam-playground-content\"\n      }\n    },\n    \"subscription\": \"projects/mikenimer-dam-playground/subscriptions/dam-gcs-fileinfo\"\n  },\n  \"resource\": {\n    \"type\": \"cloud_run_revision\",\n    \"labels\": {\n      \"service_name\": \"services-core-gcsfileinfo\",\n      \"project_id\": \"mikenimer-dam-playground\",\n      \"revision_name\": \"services-core-gcsfileinfo-00011-bud\",\n      \"configuration_name\": \"services-core-gcsfileinfo\",\n      \"location\": \"us-central1\"\n    }\n  },\n  \"timestamp\": \"2020-03-17T04:42:57.572163Z\",\n  \"labels\": {\n    \"instanceId\": \"00bf4bf02db38963e1071707cc67d40b77cc5f432ced6a5e286d27d248e3b250b574eeb5f2ee9356024259820fc6826957560318478a4b53dba7cb84a0b164329d82d8\"\n  },\n  \"logName\": \"projects/mikenimer-dam-playground/logs/run.googleapis.com%2Fstderr\",\n  \"receiveTimestamp\": \"2020-03-17T04:42:57.753990377Z\"\n}"


func TestGetFileInfo(t *testing.T) {
	name := "exif/jpg/gps/DSCN0021.jpg"
	bucket := "mikenimer-dam-playground-content"

	req, err := http.NewRequest("GET", "/?bucket=" +bucket +"&name=" +name, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RequestHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	//jsonMap := make(map[string]interface{})
	jsonAttr := storage.ObjectAttrs{}
	jsonErr := json.Unmarshal([]byte(body), &jsonAttr)
	if jsonErr != nil {
		log.Printf("json.Unmarshal: %v", jsonErr)
		return
	}

	assert := assert.New(t)
	assert.Equal(rr.Code, 200, "Incorrect Response Code")
	assert.Equal(rr.Header().Get("Content-Type"), "application/json", "Incorrect Header")

	//m, err := json.Marshal(body)
	println("Results: -----")
	println(string(body))
}

func TestGcsHandlerGetFileInfo(t *testing.T) {
	//var testMsg string = "{\"message\":{\"attributes\":{\"bucketId\":\"mikenimer-dam-playground-content\",\"eventTime\":\"2020-03-10T00:37:25.931866Z\",\"eventType\":\"OBJECT_FINALIZE\",\"notificationConfig\":\"projects/_/buckets/mikenimer-dam-playground-content/notificationConfigs/1\",\"objectGeneration\":\"1583800645932016\",\"objectId\":\"AlaskanGlacier.jpg\",\"payloadFormat\":\"JSON_API_V1\"},\"data\":\"ewogICJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwKICAiaWQiOiAibWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvQWxhc2thbkdsYWNpZXIuanBnLzE1ODM4MDA2NDU5MzIwMTYiLAogICJzZWxmTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby9BbGFza2FuR2xhY2llci5qcGciLAogICJuYW1lIjogIkFsYXNrYW5HbGFjaWVyLmpwZyIsCiAgImJ1Y2tldCI6ICJtaWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudCIsCiAgImdlbmVyYXRpb24iOiAiMTU4MzgwMDY0NTkzMjAxNiIsCiAgIm1ldGFnZW5lcmF0aW9uIjogIjEiLAogICJjb250ZW50VHlwZSI6ICJpbWFnZS9qcGVnIiwKICAidGltZUNyZWF0ZWQiOiAiMjAyMC0wMy0xMFQwMDozNzoyNS45MzFaIiwKICAidXBkYXRlZCI6ICIyMDIwLTAzLTEwVDAwOjM3OjI1LjkzMVoiLAogICJzdG9yYWdlQ2xhc3MiOiAiU1RBTkRBUkQiLAogICJ0aW1lU3RvcmFnZUNsYXNzVXBkYXRlZCI6ICIyMDIwLTAzLTEwVDAwOjM3OjI1LjkzMVoiLAogICJzaXplIjogIjcwOTYzMzIiLAogICJtZDVIYXNoIjogIlljUWM2YjNSaVozVVhGcEducFRRV2c9PSIsCiAgIm1lZGlhTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9kb3dubG9hZC9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby9BbGFza2FuR2xhY2llci5qcGc/Z2VuZXJhdGlvbj0xNTgzODAwNjQ1OTMyMDE2JmFsdD1tZWRpYSIsCiAgImNyYzMyYyI6ICJpZktWMFE9PSIsCiAgImV0YWciOiAiQ1BDbmxaZlZqdWdDRUFFPSIKfQo=\",\"messageId\":\"1037428102269218\",\"message_id\":\"1037428102269218\",\"publishTime\":\"2020-03-10T00:37:26.069Z\",\"publish_time\":\"2020-03-10T00:37:26.069Z\"},\"subscription\":\"projects/mikenimer-dam-playground/subscriptions/dam-pubsubhandler\"}"
	var testMsg string = "{\"message\":{\"attributes\":{\"bucketId\":\"mikenimer-dam-playground-content\",\"eventTime\":\"2020-03-10T00:37:25.931866Z\",\"eventType\":\"OBJECT_FINALIZE\",\"notificationConfig\":\"projects/_/buckets/mikenimer-dam-playground-content/notificationConfigs/1\",\"objectGeneration\":\"1583800645932016\",\"objectId\":\"AlaskanGlacier.jpg\",\"payloadFormat\":\"JSON_API_V1\"},\"data\":\"ewogICJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwKICAiaWQiOiAibWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvQWxhc2thbkdsYWNpZXIuanBnLzE1ODM4MDA2NDU5MzIwMTYiLAogICJzZWxmTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby9BbGFza2FuR2xhY2llci5qcGciLAogICJuYW1lIjogIkFsYXNrYW5HbGFjaWVyLmpwZyIsCiAgImJ1Y2tldCI6ICJtaWtlbmltZXItZGFtLXBsYXlncm91bmQtY29udGVudCIsCiAgImdlbmVyYXRpb24iOiAiMTU4MzgwMDY0NTkzMjAxNiIsCiAgIm1ldGFnZW5lcmF0aW9uIjogIjEiLAogICJjb250ZW50VHlwZSI6ICJpbWFnZS9qcGVnIiwKICAidGltZUNyZWF0ZWQiOiAiMjAyMC0wMy0xMFQwMDozNzoyNS45MzFaIiwKICAidXBkYXRlZCI6ICIyMDIwLTAzLTEwVDAwOjM3OjI1LjkzMVoiLAogICJzdG9yYWdlQ2xhc3MiOiAiU1RBTkRBUkQiLAogICJ0aW1lU3RvcmFnZUNsYXNzVXBkYXRlZCI6ICIyMDIwLTAzLTEwVDAwOjM3OjI1LjkzMVoiLAogICJzaXplIjogIjcwOTYzMzIiLAogICJtZDVIYXNoIjogIlljUWM2YjNSaVozVVhGcEducFRRV2c9PSIsCiAgIm1lZGlhTGluayI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9kb3dubG9hZC9zdG9yYWdlL3YxL2IvbWlrZW5pbWVyLWRhbS1wbGF5Z3JvdW5kLWNvbnRlbnQvby9BbGFza2FuR2xhY2llci5qcGc/Z2VuZXJhdGlvbj0xNTgzODAwNjQ1OTMyMDE2JmFsdD1tZWRpYSIsCiAgImNyYzMyYyI6ICJpZktWMFE9PSIsCiAgImV0YWciOiAiQ1BDbmxaZlZqdWdDRUFFPSIKfQo=\",\"messageId\":\"1037428102269218\",\"message_id\":\"1037428102269218\",\"publishTime\":\"2020-03-10T00:37:26.069Z\",\"publish_time\":\"2020-03-10T00:37:26.069Z\"},\"subscription\":\"projects/mikenimer-dam-playground/subscriptions/dam-pubsubhandler\"}\n"


	reader := strings.NewReader(testMsg)
	req, err := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RequestHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	//jsonMap := make(map[string]interface{})
	jsonAttr := storage.ObjectAttrs{}
	jsonErr := json.Unmarshal([]byte(body), &jsonAttr)
	if jsonErr != nil {
		log.Printf("json.Unmarshal: %v", jsonErr)
		return
	}

	assert := assert.New(t)
	assert.Equal(rr.Code, 200, "Incorrect Response Code")
	assert.Equal(rr.Header().Get("Content-Type"), "application/json", "Incorrect Header")

	//m, err := json.Marshal(body)
	println("Results: -----")
	println(string(body))
}
