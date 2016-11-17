package shelflib

import (
	"encoding/json"
	"github.com/tomnomnom/linkheader"
	"io"
	"net/http"
)

// Takes a response from Shelf and parses the links.
func ParseLinks(response *http.Response) (linkheader.Links, *ShelfError) {
	var (
        links linkheader.Links
        shelfErr *ShelfError
    )

	err := CheckResponseStatus(response)

	if err != nil {
        shelfErr = CreateShelfErrorFromError(err)

		return links, shelfErr
	}

	links = linkheader.ParseMultiple(response.Header["Link"])

	return links, shelfErr
}

// Parses a response with an expected JSON body.
func ParseJsonResponse(response *http.Response, result *interface{}) *ShelfError {
    var shelfErr *ShelfError

	err := CheckResponseStatus(response)

	if err != nil {
        shelfErr = CreateShelfErrorFromError(err)

		return shelfErr
	}

	loadJsonBody(response.Body, result)

	return shelfErr
}

// Parses metadata property response.
func ParseMetadataResponse(response *http.Response) (*MetadataProperty, *ShelfError) {
	var (
		jsonResponse interface{}
		result       *MetadataProperty
        shelfErr *ShelfError
	)
	err := ParseJsonResponse(response, &jsonResponse)

	if err != nil {
        shelfErr = CreateShelfErrorFromError(err)

		return result, shelfErr
	}

	prop := jsonResponse.(map[string]interface{})
	name := prop["name"].(string)
	value := prop["value"].(string)
	immutable := prop["immutable"].(bool)
	result = CreateMetadataProperty(name, value, immutable)

	return result, shelfErr
}

// Parses bulk metadata response.
func ParseBulkMetadataResponse(response *http.Response) (map[string]*MetadataProperty, *ShelfError) {
	var (
		jsonResponse interface{}
		result       map[string]*MetadataProperty
        shelfErr *ShelfError
	)

	err := ParseJsonResponse(response, &jsonResponse)

	if err != nil {
        shelfErr = CreateShelfErrorFromError(err)

		return result, shelfErr
	}

	propMap := jsonResponse.(map[string]interface{})
	result = make(map[string]*MetadataProperty)

	for key, val := range propMap {
		prop := val.(map[string]interface{})
		value := prop["value"].(string)
		immutable := prop["immutable"].(bool)
		result[key] = CreateMetadataProperty(key, value, immutable)
	}

	return result, shelfErr
}

// Creates a new MetadataProperty.
func CreateMetadataProperty(name string, value string, immutable bool) *MetadataProperty {
	mappedMetadata := &MetadataProperty{
		Name:      name,
		Value:     value,
		Immutable: immutable,
	}

	return mappedMetadata
}

// Checks given response to see if it is an error response.
// If it is it create a ShelfError.
func CheckResponseStatus(response *http.Response) *ShelfError {
    var shelfErr *ShelfError

	if response.StatusCode < 399 && response.StatusCode > 199 {
		return shelfErr
	}

	var (
		code       string
		message    string
		parsedBody interface{}
	)

	err := loadJsonBody(response.Body, &parsedBody)

	if err != nil {
        shelfErr = CreateShelfErrorFromError(err)
	} else {
		body := parsedBody.(map[string]interface{})
		message = body["message"].(string)
		code = body["code"].(string)
	    shelfErr = CreateShelfError(message, code)
	}

	return shelfErr
}

// Unamarshals JSON data from a response body.
func loadJsonBody(rawBody io.ReadCloser, result *interface{}) error {
	decoder := json.NewDecoder(rawBody)
	err := decoder.Decode(result)

	return err
}
