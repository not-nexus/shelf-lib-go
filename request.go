package shelflib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Request struct {
	Logger     log.Logger
	ShelfToken string
}

var SuffixMap = map[string]string{"meta": "_meta", "search": "_search", "artifact": ""}

// Marshals given data to JSON and creates a buffer from it.
func (this *Request) MarshalRequestData(data interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}

// Performs request on Shelf.
func (this *Request) DoRequest(verb string, path string, requestType string, property string, data io.Reader) (*http.Response, error) {
	var response http.Response
	requestURI, err := this.buildUrl(path, requestType, property)

	if err != nil {
		return response, err
	}

	req, err := http.NewRequest(verb, requestURI, data)

	if err != nil {
		return response, err
	}

	req.Header.Add("Authorization", this.ShelfToken)
	client := &http.Client{}

	return client.Do(req)
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
