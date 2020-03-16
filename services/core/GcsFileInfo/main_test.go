package main

import (
	"cloud.google.com/go/storage"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)


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
