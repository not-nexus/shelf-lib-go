package shelflib

import (
	"encoding/json"
	"fmt"
	"io"
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
func ParseJsonResponse(response *http.Response, result *interface{}) error {
	err := CheckResponseStatus(response)

	if err != nil {
		return err
	}

	loadJsonBody(response.Body, result)

	return nil
}

// Parses metadata property response.
func ParseMetadataResponse(response *http.Response) (*MetadataProperty, error) {
	var (
		jsonResponse interface{}
		result       *MetadataProperty
	)
	err := ParseJsonResponse(response, &jsonResponse)

	if err != nil {
		return result, err
	}

	prop := jsonResponse.(map[string]interface{})
	name := prop["name"].(string)
	value := prop["value"].(string)
	immutable := prop["immutable"].(bool)
	result = CreateMetadataProperty(name, value, immutable)

	return result, nil
}

// Parses bulk metadata response.
func ParseBulkMetadataResponse(response *http.Response) (map[string]*MetadataProperty, error) {
	var (
		jsonResponse interface{}
		result       map[string]*MetadataProperty
	)

	err := ParseJsonResponse(response, &jsonResponse)

	if err != nil {
		return result, err
	}

	propMap := jsonResponse.(map[string]interface{})
	result = make(map[string]*MetadataProperty)

	for key, val := range propMap {
		prop := val.(map[string]interface{})
		value := prop["value"].(string)
		immutable := prop["immutable"].(bool)
		result[key] = CreateMetadataProperty(key, value, immutable)
	}

	return result, nil
}

// Creates a new MetadataProperty.
func CreateMetadataProperty(name string, value string, immutable bool) *MetadataProperty {
	mappedMetadata := &MetadataProperty{
		Name:      name,
		Value:     value,
		Immutable: immutable,
	}

	fmt.Println(mappedMetadata)
	return mappedMetadata
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
