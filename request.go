package shelflib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

var SuffixMap = map[string]string{"meta": "_meta", "search": "_search", "artifact": ""}

func MarshalRequestData(data interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}

func DoRequest(verb string, shelfToken string, path string, requestType string, property string, data io.Reader) (*http.Response, error) {
	var response http.Response
	requestURI, err := buildUrl(path, requestType, property)

	if err != nil {
		return response, err
	}

	req, err := http.NewRequest(verb, requestURI, data)

	if err != nil {
		return response, err
	}

	req.Header.Add("Authorization", shelfToken)
	client := &http.Client{}

	return client.Do(req)
}

func buildUrl(uri string, requestType string, property string) (string, error) {
	parsedUri, err := url.Parse(uri)

	if err != nil {
		return "", err
	}

	suffix := SuffixMap[requestType]
	parsedUri.Path = path.Join(parsedUri.Path, suffix, property)

	return parsedUri.String(), nil
}
