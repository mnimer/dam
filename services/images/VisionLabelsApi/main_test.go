package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	gcp "mikenimer.com/services/core/GcpUtils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createPubSubMsg(t *testing.T, bucket, name, contentType string) string {
	// Create a request to pass to our handler.
	data := gcp.PubSubData{}
	data.Bucket = bucket
	data.Name = name
	data.ContentType = contentType

	dataJson, err := json.Marshal(data)
	if err != nil {
		assert.Fail(t, "Json Encoding Error")
	}
	msg := gcp.PubSubMessage{}
	msg.Message.ID = uuid.New().String()
	msg.Message.Data = dataJson
	msg.Subscription = base64.StdEncoding.EncodeToString(dataJson)

	testPubSubMsg, err := json.Marshal(msg)
	if err != nil {
		assert.Fail(t, "Json Encoding Error")
	}

	return string(testPubSubMsg)
}

func TestGayaSmall(t *testing.T) {
	_bucket := "mikenimer-dam-playground-content"
	_name := "unit-tests/visionapi/setagaya_small.jpeg"
	_contentType := "image/jpg"
	testPubSubMsg := createPubSubMsg(t, _bucket, _name, _contentType)
	msgBody, err := invokeObjectFinalizeMsg(t, string(testPubSubMsg))

	if err != nil {
		log.Printf(err.Error())
		t.Fail()
	} else {
		assert := assert.New(t)
		assert.Equal(_name, msgBody.Name)
		assert.Equal(_bucket, msgBody.Bucket)
		assert.GreaterOrEqual(len(msgBody.Metadata), 10)
	}
}

func invokeObjectFinalizeMsg(t *testing.T, msg string) (LabelsApiResult, error) {
	//file := "https://www.googleapis.com/storage/v1/b/mikenimer-dam-playground-content/o/AlaskanGlacier.jpg";

	reader := strings.NewReader(msg)
	req, err := http.NewRequest("GET", "/", reader)
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

	var msgBody LabelsApiResult

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		return msgBody, errors.New(rr.Body.String())
	}
	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &msgBody); err != nil {
		//log.Fatal("json.Unmarshal: " +err.Error())
		return msgBody, errors.New(err.Error())
	}

	return msgBody, nil
}
