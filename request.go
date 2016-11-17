package shelflib

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

type Request struct {
	Logger     *log.Logger
	ShelfToken string
}

var SuffixMap = map[string]string{"meta": "_meta", "search": "_search", "artifact": ""}

// Marshals given data to JSON and creates a buffer from it.
func (this *Request) MarshalRequestData(data interface{}) (io.Reader, *ShelfError) {
	var shelfErr *ShelfError

	jsonData, err := json.Marshal(data)

	if err != nil {
		shelfErr = CreateShelfErrorFromError(err)

		return nil, shelfErr
	}

	return bytes.NewBuffer(jsonData), shelfErr
}

// Performs request on Shelf.
func (this *Request) DoRequest(verb string, path string, requestType string, property string, data io.Reader) (*http.Response, *ShelfError) {
	var shelfErr *ShelfError

	requestURI, err := this.buildUrl(path, requestType, property)

	if err != nil {
		shelfErr = CreateShelfErrorFromError(err)

		return nil, shelfErr
	}

	req, err := http.NewRequest(verb, requestURI, data)

	if err != nil {
		shelfErr = CreateShelfErrorFromError(err)

		return nil, shelfErr
	}

	req.Header.Add("Authorization", this.ShelfToken)
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		shelfErr = CreateShelfErrorFromError(err)

		return nil, shelfErr
	}

	return resp, shelfErr
}

// Builds Shelf URL.
func (this *Request) buildUrl(uri string, requestType string, property string) (string, error) {
	parsedUri, err := url.Parse(uri)

	if err != nil {
		return "", err
	}

	suffix := SuffixMap[requestType]
	parsedUri.Path = path.Join(parsedUri.Path, suffix, property)

	return parsedUri.String(), nil
}
