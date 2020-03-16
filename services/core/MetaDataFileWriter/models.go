package main

type Metadata struct {
	FileId string
	Bucket string
	Name string
	Metadata map[string]interface{}
}

type Geo struct {
	Latitude float64
	Longitude float64
}
