package Gcp

import (
	"io/ioutil"
	"net/http"
)

func GetMetadataFolder(fileName string) string {
	var dirName = ".metadata/" +fileName
	return dirName
}



func GetProjectId() (string, error) {
	url := "http://metadata.google.internal/computeMetadata/v1/project/project-id"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := client.Do(req)

	if err != nil {
		println("Error:" +err.Error())
		return "mikenimer-dam-playground", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Error:" +err.Error())
		return "mikenimer-dam-playground", err
	}
	println("Project: " +string(body))
	return string(body), nil
}
