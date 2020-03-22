package main

type ExifPrimary struct {
	Bucket	string
	Name 	string
	Metadata PrimaryTags
}
type ExifGps struct {
	Bucket 	string
	Name 	string
	Metadata GpsTags
}
type ExifGeo struct {
	Bucket 	string
	Name 	string
	Metadata Geo
}
type ExifTags struct {
	Bucket 	string
	Name 	string
	Metadata map[string]interface{}
}



type Geo struct {
	Latitude   	float64	`json:"latitude"`
	Longitude 	float64	`json:"latitude"`
}

type PrimaryTags struct {
	XPTitle    string	`json:"xpTitle,omitempty"`
	XPComment  string	`json:"xpComment,omitempty"`
	XPAuthor   string	`json:"xpAuthor,omitempty"`
	XPKeywords string	`json:"xpKeywords,omitempty"`
	XPSubject  string	`json:"xpSubject,omitempty"`
}

type GpsTags struct {
	GPSVersionID        string	`json:"gpsVersionId,omitempty"`
	GPSLatitudeRef      string	`json:"gpsLatitudeRef,omitempty"`
	GPSLatitude         float64	`json:"gpsLatitude,omitempty"`
	GPSLongitudeRef     string	`json:"gpsLongitudeRef,omitempty"`
	GPSLongitude        float64	`json:"gpsLongitude,omitempty"`
	GPSAltitudeRef      string	`json:"gpsAltitudeRef,omitempty"`
	GPSAltitude         string	`json:"gpsAltitude,omitempty"`
	GPSTimeStamp        string	`json:"gpsTimeStamp,omitempty"`
	GPSSatelites        string	`json:"gpsSatelites,omitempty"`
	GPSStatus           string	`json:"gpsStatus,omitempty"`
	GPSMeasureMode      string	`json:"gpsMeasureMode,omitempty"`
	GPSDOP              string	`json:"gpsDop,omitempty"`
	GPSSpeedRef         string	`json:"gpsSpeedRef,omitempty"`
	GPSSpeed            string	`json:"gpsSpeed,omitempty"`
	GPSTrackRef         string	`json:"gpsTrackRef,omitempty"`
	GPSTrack            string	`json:"gpsTrack,omitempty"`
	GPSImgDirectionRef  string	`json:"gpsImgDirectionRef,omitempty"`
	GPSImgDirection     string	`json:"gpsImgDirection,omitempty"`
	GPSMapDatum         string	`json:"gpsMapDatum,omitempty"`
	GPSDestLatitudeRef  string	`json:"gpsDestLatitudeRef,omitempty"`
	GPSDestLatitude     string	`json:"gpsDestLatitude,omitempty"`
	GPSDestLongitudeRef string	`json:"gpsDestLongitudeRef,omitempty"`
	GPSDestLongitude    string	`json:"gpsDestLongitude,omitempty"`
	GPSDestBearingRef   string	`json:"gpsDestBearingRef,omitempty"`
	GPSDestBearing      string	`json:"gpsDestBearing,omitempty"`
	GPSDestDistanceRef  string	`json:"gpsDestDistanceRef,omitempty"`
	GPSDestDistance     string	`json:"gpsDestDistance,omitempty"`
	GPSProcessingMethod string	`json:"gpsProcessingMethod,omitempty"`
	GPSAreaInformation  string	`json:"gpsAreaInformation,omitempty"`
	GPSDateStamp        string	`json:"gpsDateStamp,omitempty"`
	GPSDifferential     string	`json:"gpsDifferential,omitempty"`
}
