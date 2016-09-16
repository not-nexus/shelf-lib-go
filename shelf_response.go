package shelflib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Parses a response from Shelf.
func ParseShelfResponse(response *http.Response) (io.ReadCloser, error) {
	var err error

	if response.StatusCode > 399 {
		err := checkStatus(response.Body)
	}

	return data, err
}

// Creates an error from a Shelf error response.
func checkStatus(body io.ReadCloser) error {
	var (
		message string
		code    string
	)
	body, err := loadJsonBody(body)

	if err != nil {
		message = ""
		code = ""
	} else {
		message = body.code
		code = body.message
	}

	err = NewShelfError(message, code)

	return err
}

// Unamarshals JSON data from a response body.
func loadJsonBody(rawBody io.ReadCloser) (interface{}, error) {
	body, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	var resp interface{}
	err = json.Unmarshal(body, &resp)

	return resp, err
}
