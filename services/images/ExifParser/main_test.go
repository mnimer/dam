package main

import (
	"encoding/base64"
	"encoding/json"
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


func createPubSubMsg(t *testing.T, bucket, name, contentType string) string {
	// Create a request to pass to our handler.
	data := gcp.PubSubData{}
	data.Bucket = bucket
	data.Name = name
	data.ContentType = contentType

	dataJson, err := json.Marshal(data);
	if( err != nil){
		assert.Fail(t, "Json Encoding Error")
	}
	msg := gcp.PubSubMessage{}
	msg.Message.ID = uuid.New().String()
	msg.Message.Data = dataJson;
	msg.Subscription = base64.StdEncoding.EncodeToString(dataJson)

	testPubSubMsg, err := json.Marshal(msg);
	if( err != nil){
		assert.Fail(t, "Json Encoding Error")
	}

	return string(testPubSubMsg)
}


func TestJpgMsg(t *testing.T) {
	_bucket := "mikenimer-dam-playground-content"
	_name := "unit-tests/exif/jpg/Canon_40D.jpg"
	_contentType := "image/jpg"
	testPubSubMsg := createPubSubMsg(t, _bucket, _name, _contentType)
	msgBody := invokeObjectFinalizeMsg(t, string(testPubSubMsg))

	assert := assert.New(t)
	assert.Equal(msgBody["Name"], _name)
	assert.Equal(msgBody["Bucket"], _bucket)

	//Test random metadata KV to make sure the file was parsed
	var md map[string]interface{} = msgBody["Metadata"].(map[string]interface{});
	assert.Equal(md["Technology"], "Cathode Ray Tube Display")
	assert.Equal(md["DeviceManufacturer"], "Hewlett-Packard")
	assert.Equal(md["Make"], "Canon")
	//make sure these pointer to our tmp file are not included
	assert.Nil(md["Directory"])
	assert.Nil(md["FileName"])
	assert.Nil(md["SourceFile"])
}


func TestJpgWithGPSMsg(t *testing.T) {
	_bucket := "mikenimer-dam-playground-content"
	_name := "unit-tests/exif/jpg/gps/DSCN0010.jpg"
	_contentType := "image/jpg"
	testPubSubMsg := createPubSubMsg(t, _bucket, _name, _contentType)
	msgBody := invokeObjectFinalizeMsg(t, string(testPubSubMsg))

	assert := assert.New(t)
	assert.Equal(msgBody["Name"], _name)
	assert.Equal(msgBody["Bucket"], _bucket)

	//Test random metadata KV to make sure the file was parsed
	var md map[string]interface{} = msgBody["Metadata"].(map[string]interface{});
	assert.Equal("COOLPIX P6000", md["Model"], )
	assert.Equal("06", md["GPSSatellites"])
	assert.Equal("43 deg 28' 2.81\" N", md["GPSLatitude"])
	assert.Equal("11 deg 53' 6.46\" E", md["GPSLongitude"])
	//make sure these pointer to our tmp file are not included
	assert.Nil(md["Directory"])
	assert.Nil(md["FileName"])
	assert.Nil(md["SourceFile"])
}

func TestTiffMsg(t *testing.T) {
	_bucket := "mikenimer-dam-playground-content"
	_name := "unit-tests/exif/tiff/Arbitro.tiff"
	_contentType := "image/jpg"
	testPubSubMsg := createPubSubMsg(t, _bucket, _name, _contentType)
	msgBody := invokeObjectFinalizeMsg(t, string(testPubSubMsg))

	assert := assert.New(t)
	assert.Equal(msgBody["Name"], _name)
	assert.Equal(msgBody["Bucket"], _bucket)

	//Test random metadata KV to make sure the file was parsed
	var md map[string]interface{} = msgBody["Metadata"].(map[string]interface{});
	assert.Equal("TIFF", md["FileType"] )
	assert.Equal(0.007, md["Megapixels"])
	assert.Equal("174x38", md["ImageSize"] )
	assert.Equal("RGB", md["PhotometricInterpretation"])
	//make sure these pointer to our tmp file are not included
	assert.Nil(md["Directory"])
	assert.Nil(md["FileName"])
	assert.Nil(md["SourceFile"])
}


func invokeObjectFinalizeMsg(t *testing.T, msg string) map[string]interface{} {
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

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var msgBody map[string]interface{}
	//Parse the rest results
	body, err := ioutil.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &msgBody); err != nil {
		log.Fatal("json.Unmarshal: " +err.Error())
	}else {
		assert := assert.New(t)
		assert.Equal(rr.Code, 200, "Incorrect Response Code")
	}

	return msgBody
}


