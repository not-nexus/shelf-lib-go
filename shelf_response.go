package shelflib

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// Takes a response from Shelf and parses the links.
func ParseLinks(response *http.Response) ([]string, error) {
	var links []string

	err := CheckResponseStatus(response)

	if err != nil {
		return links, err
	}

	return response.Header["Links"], nil
}

// Parses a response with an expected JSON body.
func ParseJsonResponse(response *http.Response, result interface{}) error {
	err := CheckResponseStatus(response)

	if err != nil {
		return err
	}

	loadJsonBody(response.Body, &result)

	return nil
}

func ParseStreamResponse(response *http.Response) ([]byte, error) {
	err := CheckResponseStatus(response)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}

// Checks given response to see if it is an error response.
// If it is it create a ShelfError.
func CheckResponseStatus(response *http.Response) error {
	if response.StatusCode < 399 {
		return nil
	}

	var (
		code       string
		message    string
		parsedBody interface{}
	)

	err := loadJsonBody(response.Body, &parsedBody)

	if err != nil {
		return err
	} else {
		body := parsedBody.(map[string]interface{})
		message = body["message"].(string)
		code = body["code"].(string)
	}

	shelfErr := NewShelfError(message, code)

	return shelfErr
}

// Unamarshals JSON data from a response body.
func loadJsonBody(rawBody io.ReadCloser, result *interface{}) error {
	decoder := json.NewDecoder(rawBody)
	err := decoder.Decode(result)

	return err
}
