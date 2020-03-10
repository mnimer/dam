package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestObjectFinalizeMsg(t *testing.T) {
	// Create a request to pass to our handler.
	name := "AlaskanGlacier.jpg";
	bucket := "mikenimer-dam-playground-content"
	file := "mikenimer-dam-playground-content/AlaskanGlacier.jpg/1583800645932016";
	//file := "https://www.googleapis.com/storage/v1/b/mikenimer-dam-playground-content/o/AlaskanGlacier.jpg";
	req, err := http.NewRequest("GET", "/?name=" +name +"&bucket=" +bucket +"&file=" +file, nil)
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

	var msgBody map[string]interface{}
	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &msgBody); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		return
	}

	assert := assert.New(t)
	assert.Equal(rr.Code, 200, "Incorrect Response Code")
	assert.Equal(rr.Header().Get("Content-Type"), "application/json", "Incorrect Header")
}