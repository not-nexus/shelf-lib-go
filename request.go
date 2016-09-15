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

// ShelfResponse is a wrapper for a response from shelf.
type ShelfResponse struct {
	Body       interface{}
	Links      []string
	StatusCode int
}

func MarshalRequestData(data interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}

func DoRequest(verb string, shelfToken string, path string, requestType string, property string, data io.Reader) (*ShelfResponse, error) {
	response := &ShelfResponse{}
	requestURI, err := buildUrl(path, requestType, property)

	if err != nil {
		return response, err
	}

	req, err := http.NewRequest(verb, requestURI, data)

	if err != nil {
		return response, err
	}

	req.Header.Add("Authorization", shelfReq.ShelfToken)
	client := &http.Client{}
	rawResponse, err := client.Do(req)

	if err != nil {
		return response, err
	}

	return parseRequest(rawResponse), nil
}

func parseRequest(response *http.Response) *ShelfResponse {
	shelfResponse := &ShelfResponse{
		StatusCode: response.StatusCode,
		Links:      response.Header["Link"],
		Body:       response.Body,
	}

	return shelfResponse
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
