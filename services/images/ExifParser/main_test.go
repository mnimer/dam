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

func TestParseCordLat(t *testing.T) {
	lat := parseCoordString("43 deg 28' 2.81\" N")
	assert.Equal(t, 43.467447222222226, lat, "Latitude is wrong")
}
func TestParseCordLat2(t *testing.T) {
	lat := parseCoordString("43 deg 28' 2.81\" S")
	assert.Equal(t, -43.467447222222226, lat, "Latitude is wrong")
}


func TestParseCordLng(t *testing.T) {
	lng := parseCoordString("11 deg 53' 6.46\" E")
	assert.Equal(t, 11.885127777777777, lng, "Longitude is wrong")
}
func TestParseCordLng2(t *testing.T) {
	lng := parseCoordString("11 deg 53' 6.46\" W")
	assert.Equal(t, -11.885127777777777, lng, "Longitude is wrong")
}


func TestJpgWithGPSMsg(t *testing.T) {
	// Create a request to pass to our handler.
	name := "exif/jpg/gps/DSCN0021.jpg"
	bucket := "mikenimer-dam-playground-content"
	file := bucket +"/" +name;
	runObjectFinalizeMsg(t, name, bucket, file)
}



func TestTiffWithGPSMsg(t *testing.T) {
	// Create a request to pass to our handler.
	name := "exif/tiff/Arbitro.tiff"
	bucket := "mikenimer-dam-playground-content"
	file := bucket +"/" +name;
	runObjectFinalizeMsg(t, name, bucket, file)
}

func runObjectFinalizeMsg(t *testing.T, name, bucket, file string) {
	//file := "https://www.googleapis.com/storage/v1/b/mikenimer-dam-playground-content/o/AlaskanGlacier.jpg";
	req, err := http.NewRequest("GET", "/?name=" +name +"&bucket=" +bucket +"&file=" +file, nil)
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

	var msgBody Exif
	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &msgBody); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		return
	}

	assert := assert.New(t)
	assert.Equal(rr.Code, 200, "Incorrect Response Code")
	assert.Equal(rr.Header().Get("Content-Type"), "application/json", "Incorrect Header")
	assert.Equal(msgBody.Name, name)
	assert.Equal(msgBody.Bucket, bucket)
	assert.Equal(msgBody.FileId, file)

	m, err := json.Marshal(msgBody)
	println("Results: -----")
	println(string(m))
}


