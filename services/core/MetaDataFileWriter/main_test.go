package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetadataWriter(t *testing.T) {
	input := "{\"FileId\":\"mikenimer-dam-playground-content/exif/jpg/gps/DSCN0010.jpg\",\"Bucket\":\"mikenimer-dam-playground-content\",\"Name\":\"exif/jpg/gps/DSCN0021.jpg\",\"Metadata\":{\"PrimaryTags\":{\"XPTitle\":\"\",\"XPComment\":\"\",\"XPAuthor\":\"\",\"XPKeywords\":\"\",\"XPSubject\":\"\"},\"GpsTags\":{\"GPSVersionID\":\"\",\"GPSLatitudeRef\":\"North\",\"GPSLatitude\":\"43.467081\",\"GPSLongitudeRef\":\"East\",\"GPSLongitude\":\"11.884539\",\"GPSAltitudeRef\":\"Above Sea Level\",\"GPSAltitude\":\"\",\"GPSTimeStamp\":\"14:36:47.23\",\"GPSSatelites\":\"\",\"GPSStatus\":\"\",\"GPSMeasureMode\":\"\",\"GPSDOP\":\"\",\"GPSSpeedRef\":\"\",\"GPSSpeed\":\"\",\"GPSTrackRef\":\"\",\"GPSTrack\":\"\",\"GPSImgDirectionRef\":\"Unknown ()\",\"GPSImgDirection\":\"\",\"GPSMapDatum\":\"WGS-84   \",\"GPSDestLatitudeRef\":\"\",\"GPSDestLatitude\":\"\",\"GPSDestLongitudeRef\":\"\",\"GPSDestLongitude\":\"\",\"GPSDestBearingRef\":\"\",\"GPSDestBearing\":\"\",\"GPSDestDistanceRef\":\"\",\"GPSDestDistance\":\"\",\"GPSProcessingMethod\":\"\",\"GPSAreaInformation\":\"\",\"GPSDateStamp\":\"2008:10:23\",\"GPSDifferential\":\"\"},\"JpegTags\":{\"ImageWidth\":\"\",\"ImageLength\":\"\",\"BitsPerSample\":\"\",\"Compression\":\"JPEG (old-style)\",\"PhotometricInterpretation\":\"\",\"Orientation\":\"Horizontal (normal)\",\"SamplesPerPixel\":\"\",\"PlanarConfiguration\":\"\",\"YCbCrSubSampling\":\"YCbCr4:2:2 (2 1)\",\"YCbCrPositioning\":\"Centered\",\"XResolution\":\"\",\"YResolution\":\"\",\"ResolutionUnit\":\"inches\",\"DateTime\":\"\",\"ImageDescription\":\"                               \",\"Make\":\"NIKON\",\"Model\":\"COOLPIX P6000\",\"Software\":\"Nikon Transfer 1.1 W\",\"Artist\":\"\",\"Copyright\":\"\",\"ExifIFDPointer\":\"\",\"GPSInfoIFDPointer\":\"\",\"InteroperabilityIFDPointer\":\"\",\"ExifVersion\":\"0220\",\"FlashpixVersion\":\"0100\",\"ColorSpace\":\"sRGB\",\"ComponentsConfiguration\":\"Y, Cb, Cr, -\",\"CompressedBitsPerPixel\":\"\",\"PixelXDimension\":\"\",\"PixelYDimension\":\"\",\"MakerNote\":\"\",\"UserComment\":\"\",\"RelatedSoundFile\":\"\",\"DateTimeOriginal\":\"2008:10:22 16:38:20\",\"DateTimeDigitized\":\"\",\"SubSecTime\":\"\",\"SubSecTimeOriginal\":\"\",\"SubSecTimeDigitized\":\"\",\"ImageUniqueID\":\"\",\"ExposureTime\":\"1/96\",\"FNumber\":\"\",\"ExposureProgram\":\"Program AE\",\"SpectralSensitivity\":\"\",\"ISOSpeedRatings\":\"\",\"OECF\":\"\",\"ShutterSpeedValue\":\"\",\"ApertureValue\":\"\",\"BrightnessValue\":\"\",\"ExposureBiasValue\":\"\",\"MaxApertureValue\":\"\",\"SubjectDistance\":\"\",\"MeteringMode\":\"Multi-segment\",\"LightSource\":\"Unknown\",\"Flash\":\"Off, Did not fire\",\"FocalLength\":\"16.6 mm\",\"SubjectArea\":\"\",\"FlashEnergy\":\"\",\"SpatialFrequencyResponse\":\"\",\"FocalPlaneXResolution\":\"\",\"FocalPlaneYResolution\":\"\",\"FocalPlaneResolutionUnit\":\"\",\"SubjectLocation\":\"\",\"ExposureIndex\":\"\",\"SensingMethod\":\"\",\"FileSource\":\"Digital Camera\",\"SceneType\":\"Directly photographed\",\"CFAPattern\":\"\",\"CustomRendered\":\"Normal\",\"ExposureMode\":\"Auto\",\"WhiteBalance\":\"Auto\",\"DigitalZoomRatio\":\"\",\"FocalLengthIn35mmFilm\":\"\",\"SceneCaptureType\":\"Standard\",\"GainControl\":\"None\",\"Contrast\":\"Normal\",\"Saturation\":\"Normal\",\"Sharpness\":\"Normal\",\"DeviceSettingDescription\":\"\",\"SubjectDistanceRange\":\"Unknown\",\"LensMake\":\"\",\"LensModel\":\"\"}}}\n"

	reader := strings.NewReader(input)
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

}


func TestWithMissingData(t *testing.T) {
	input := "{\"FileId\":\"mikenimer-dam-playground-content/exif/jpg/gps/DSCN0010.jpg\",\"Metadata\":{\"PrimaryTags\":{\"XPTitle\":\"\",\"XPComment\":\"\",\"XPAuthor\":\"\",\"XPKeywords\":\"\",\"XPSubject\":\"\"},\"GpsTags\":{\"GPSVersionID\":\"\",\"GPSLatitudeRef\":\"North\",\"GPSLatitude\":\"43.467081\",\"GPSLongitudeRef\":\"East\",\"GPSLongitude\":\"11.884539\",\"GPSAltitudeRef\":\"Above Sea Level\",\"GPSAltitude\":\"\",\"GPSTimeStamp\":\"14:36:47.23\",\"GPSSatelites\":\"\",\"GPSStatus\":\"\",\"GPSMeasureMode\":\"\",\"GPSDOP\":\"\",\"GPSSpeedRef\":\"\",\"GPSSpeed\":\"\",\"GPSTrackRef\":\"\",\"GPSTrack\":\"\",\"GPSImgDirectionRef\":\"Unknown ()\",\"GPSImgDirection\":\"\",\"GPSMapDatum\":\"WGS-84   \",\"GPSDestLatitudeRef\":\"\",\"GPSDestLatitude\":\"\",\"GPSDestLongitudeRef\":\"\",\"GPSDestLongitude\":\"\",\"GPSDestBearingRef\":\"\",\"GPSDestBearing\":\"\",\"GPSDestDistanceRef\":\"\",\"GPSDestDistance\":\"\",\"GPSProcessingMethod\":\"\",\"GPSAreaInformation\":\"\",\"GPSDateStamp\":\"2008:10:23\",\"GPSDifferential\":\"\"},\"JpegTags\":{\"ImageWidth\":\"\",\"ImageLength\":\"\",\"BitsPerSample\":\"\",\"Compression\":\"JPEG (old-style)\",\"PhotometricInterpretation\":\"\",\"Orientation\":\"Horizontal (normal)\",\"SamplesPerPixel\":\"\",\"PlanarConfiguration\":\"\",\"YCbCrSubSampling\":\"YCbCr4:2:2 (2 1)\",\"YCbCrPositioning\":\"Centered\",\"XResolution\":\"\",\"YResolution\":\"\",\"ResolutionUnit\":\"inches\",\"DateTime\":\"\",\"ImageDescription\":\"                               \",\"Make\":\"NIKON\",\"Model\":\"COOLPIX P6000\",\"Software\":\"Nikon Transfer 1.1 W\",\"Artist\":\"\",\"Copyright\":\"\",\"ExifIFDPointer\":\"\",\"GPSInfoIFDPointer\":\"\",\"InteroperabilityIFDPointer\":\"\",\"ExifVersion\":\"0220\",\"FlashpixVersion\":\"0100\",\"ColorSpace\":\"sRGB\",\"ComponentsConfiguration\":\"Y, Cb, Cr, -\",\"CompressedBitsPerPixel\":\"\",\"PixelXDimension\":\"\",\"PixelYDimension\":\"\",\"MakerNote\":\"\",\"UserComment\":\"\",\"RelatedSoundFile\":\"\",\"DateTimeOriginal\":\"2008:10:22 16:38:20\",\"DateTimeDigitized\":\"\",\"SubSecTime\":\"\",\"SubSecTimeOriginal\":\"\",\"SubSecTimeDigitized\":\"\",\"ImageUniqueID\":\"\",\"ExposureTime\":\"1/96\",\"FNumber\":\"\",\"ExposureProgram\":\"Program AE\",\"SpectralSensitivity\":\"\",\"ISOSpeedRatings\":\"\",\"OECF\":\"\",\"ShutterSpeedValue\":\"\",\"ApertureValue\":\"\",\"BrightnessValue\":\"\",\"ExposureBiasValue\":\"\",\"MaxApertureValue\":\"\",\"SubjectDistance\":\"\",\"MeteringMode\":\"Multi-segment\",\"LightSource\":\"Unknown\",\"Flash\":\"Off, Did not fire\",\"FocalLength\":\"16.6 mm\",\"SubjectArea\":\"\",\"FlashEnergy\":\"\",\"SpatialFrequencyResponse\":\"\",\"FocalPlaneXResolution\":\"\",\"FocalPlaneYResolution\":\"\",\"FocalPlaneResolutionUnit\":\"\",\"SubjectLocation\":\"\",\"ExposureIndex\":\"\",\"SensingMethod\":\"\",\"FileSource\":\"Digital Camera\",\"SceneType\":\"Directly photographed\",\"CFAPattern\":\"\",\"CustomRendered\":\"Normal\",\"ExposureMode\":\"Auto\",\"WhiteBalance\":\"Auto\",\"DigitalZoomRatio\":\"\",\"FocalLengthIn35mmFilm\":\"\",\"SceneCaptureType\":\"Standard\",\"GainControl\":\"None\",\"Contrast\":\"Normal\",\"Saturation\":\"Normal\",\"Sharpness\":\"Normal\",\"DeviceSettingDescription\":\"\",\"SubjectDistanceRange\":\"Unknown\",\"LensMake\":\"\",\"LensModel\":\"\"}}}\n"

	reader := strings.NewReader(input)
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

}
