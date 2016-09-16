package shelflib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseShelfResponse(response *http.Response, data interface{}) (interface{}, error) {
	if response.StatusCode > 399 {
		err := checkStatus(response.Body)
	} else {

	}

	return data, err
}

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

func loadJsonBody(rawBody io.ReadCloser) (interface{}, error) {
	body, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	var resp interface{}
	err = json.Unmarshal(body, &resp)

	return resp, err
}
