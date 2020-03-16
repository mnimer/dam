package main

type Exif struct {
	FileId string
	Bucket string
	Name string
	Metadata ExifMetadata
}

type ExifMetadata struct {
	PrimaryTags PrimaryTags
	GpsTags     GpsTags
	ExifTags    map[string]interface{}
}

type PrimaryTags struct {
	XPTitle    string
	XPComment  string
	XPAuthor   string
	XPKeywords string
	XPSubject  string
}

type GpsTags struct {
	GPSVersionID        string
	GPSLatitudeRef      string
	GPSLatitude         string
	GPSLongitudeRef     string
	GPSLongitude        string
	GPSAltitudeRef      string
	GPSAltitude         string
	GPSTimeStamp        string
	GPSSatelites        string
	GPSStatus           string
	GPSMeasureMode      string
	GPSDOP              string
	GPSSpeedRef         string
	GPSSpeed            string
	GPSTrackRef         string
	GPSTrack            string
	GPSImgDirectionRef  string
	GPSImgDirection     string
	GPSMapDatum         string
	GPSDestLatitudeRef  string
	GPSDestLatitude     string
	GPSDestLongitudeRef string
	GPSDestLongitude    string
	GPSDestBearingRef   string
	GPSDestBearing      string
	GPSDestDistanceRef  string
	GPSDestDistance     string
	GPSProcessingMethod string
	GPSAreaInformation  string
	GPSDateStamp        string
	GPSDifferential     string
}
