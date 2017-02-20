package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func downloadImage(url string) (string, error) {
	filename := ""
	response, err := http.Get(url)
	if err != nil {
		return filename, err
	}

	// response.Status
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return filename, err
	}

	_, filename = path.Split(url)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return filename, err
	}

	defer func() {
		file.Close()
	}()

	file.Write(body)
	return filename, nil
}
