package shelflib

import (
	"bytes"
	"log"
)

// Wrapper for Shelf search criteria.
type SearchCriteria struct {
	Search []string
	Sort   []string
	Limit  int
}

// Wrapper for Shelf metadata property.
type MetadataProperty struct {
	name      string
	value     string
	immutable bool
}

// Interface for interacting with Shelf.
type ShelfLib struct {
	Logger  log.Logger
	Request *Request
}

// Create a ShelfLib instance.
func New(shelfToken string, logger log.Logger) *ShelfLib {
	request := &Request{Logger: logger, ShelfToken: shelfToken}

	return &ShelfLib{Logger: logger, Request: request}
}

// Download artifact from Shelf.
func (this *ShelfLib) GetArtifact(path string) ([]byte, error) {
	var artifact []byte
	response, err := this.Request.DoRequest("GET", path, "artifact", "", nil)

	if err != nil {
		return artifact, err
	}

	return ParseStreamResponse(response)
}

// Perform a HEAD request on an artifact endpoint.
func (this *ShelfLib) ListArtifact(path string) ([]string, error) {
	var links []string
	response, err := this.Request.DoRequest("HEAD", path, "artifact", "", nil)

	if err != nil {
		return links, err
	}

	return ParseLinks(response)
}

// Upload an artifact from Shelf.
func (this *ShelfLib) CreateArtifact(path string, data []byte) error {
	response, err := this.Request.DoRequest("POST", path, "artifact", "", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	return CheckResponseStatus(response)
}

// Search Shelf using SearchCriteria wrapper struct.
func (this *ShelfLib) Search(path string, searchCriteria *SearchCriteria) ([]string, error) {
	var links []string
	data, err := this.Request.MarshalRequestData(searchCriteria)

	if err != nil {
		return links, err
	}

	response, err := this.Request.DoRequest("POST", path, "search", "", data)

	if err != nil {
		return links, err
	}

	return ParseLinks(response)
}

// Retrieve metadata for an artifact.
func (this *ShelfLib) GetMetadata(path string) (map[string]*MetadataProperty, error) {
	var responseMeta map[string]*MetadataProperty
	response, err := this.Request.DoRequest("GET", path, "meta", "", nil)

	if err != nil {
		return responseMeta, nil
	}

	err = ParseJsonResponse(response, responseMeta)

	return responseMeta, err
}

// Retrieve metadata property for an artifact.
func (this *ShelfLib) GetMetadataProperty(path string, propertyKey string) (*MetadataProperty, error) {
	var responseMeta *MetadataProperty
	response, err := this.Request.DoRequest("GET", path, "meta", propertyKey, nil)

	if err != nil {
		return responseMeta, err
	}

	err = ParseJsonResponse(response, responseMeta)

	return responseMeta, err
}

// Bulk update of an artifacts metadata.
func (this *ShelfLib) UpdateMetadata(path string, metadata map[string]*MetadataProperty) (map[string]*MetadataProperty, error) {
	var responseMeta map[string]*MetadataProperty
	data, err := this.Request.MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := this.Request.DoRequest("PUT", path, "meta", "", data)

	if err != nil {
		return responseMeta, err
	}

	err = ParseJsonResponse(response, responseMeta)

	return responseMeta, err
}

// Update metadata property for an artifact.
func (this *ShelfLib) UpdateMetadataProperty(path string, metadata *MetadataProperty) (*MetadataProperty, error) {
	var responseMeta *MetadataProperty

	data, err := this.Request.MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := this.Request.DoRequest("PUT", path, "meta", metadata.name, data)

	if err != nil {
		return responseMeta, err
	}

	err = ParseJsonResponse(response, responseMeta)

	return responseMeta, err
}

// Create metadata property. Will not update existing.
func (this *ShelfLib) CreateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	var responseMeta *MetadataProperty
	data, err := this.Request.MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := this.Request.DoRequest("POST", path, "meta", metadata.name, data)

	if err != nil {
		return responseMeta, err
	}

	err = ParseJsonResponse(response, responseMeta)

	return responseMeta, err
}
