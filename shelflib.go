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
	Name      string
	Value     string
	Immutable bool
}

// Interface for interacting with Shelf.
type ShelfLib struct {
	Logger  log.Logger
	Request Request
}

// Create a ShelfLib instance.
func New(shelfToken string, logger log.Logger) *ShelfLib {
	request := &Request{logger, shelfToken}

	return &ShelfLib{logger, request}
}

// Download artifact from Shelf.
func (this *ShelfLib) GetArtifact(path string) ([]byte, error) {
	var artifact []byte
	response, err := this.Request.DoRequest("GET", path, "artifact", "", nil)

	if err != nil {
		return artifact, err
	}

	err = ParseShelfResponse(response)

	return resp.([]byte), err
}

// Perform a HEAD request on an artifact endpoint.
func (this *ShelfLib) ListArtifact(path string) ([]string, error) {
	var links []string
	response, err := this.Request.DoRequest("HEAD", path, "artifact", "", nil)

	if err != nil {
		return links, err
	}

	resp, err := ParseShelfResponse(response)

	return resp.([]string), err
}

// Upload an artifact from Shelf.
func (this *ShelfLib) CreateArtifact(path string, data []byte) error {
	response, err := this.Request.DoRequest("POST", path, "artifact", "", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	_, err := ParseShelfResponse(response)

	return err
}

// Search Shelf using SearchCriteria wrapper struct.
func (this *ShelfLib) Search(path string, searchCriteria *SearchCriteria) ([]string, error) {
	var links []string
	data, err := MarshalRequestData(searchCriteria)

	if err != nil {
		return links, err
	}

	response, err := this.Request.DoRequest("POST", path, "search", "", data)

	if err != nil {
		return links, err
	}

	resp, err := ParseShelfResponse(response)

	return resp.([]string), err
}

// Retrieve metadata for an artifact.
func (this *ShelfLib) GetMetadata(path string) (map[string]MetadataProperty, error) {
	var responseMeta map[string]*MetadataProperty
	response, err := this.Request.DoRequest("GET", path, "meta", "", nil)

	if err != nil {
		return responseMeta, nil
	}

	resp, err := ParseShelfResponse(response)

	return resp.(map[string]MetadataProperty), err
}

// Retrieve metadata property for an artifact.
func (this *ShelfLib) GetMetadataProperty(path string, propertyKey string) (*MetadataProperty, error) {
	var responseMeta MetadataProperty
	response, err := this.Request.DoRequest("GET", path, "meta", propertyKey, nil)

	if err != nil {
		return responseMeta, err
	}

	resp, err := ParseShelfResponse(response)

	return resp.(MetadataProperty), err
}

// Bulk update of an artifacts metadata.
func (this *ShelfLib) UpdateMetadata(path string, metadata map[string]*MetadataProperty) (map[string]*MetadataProperty, error) {
	var responseMeta map[string]*MetadataProperty
	data, err := MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := this.Request.DoRequest("PUT", path, "meta", "", data)

	if err != nil {
		return responseMeta, err
	}

	return ParseShelfResponse(response)
}

// Update metadata property for an artifact.
func (this *ShelfLib) UpdateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	var responseMeta MetadataProperty

	data, err := MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := this.Request.DoRequest("PUT", path, "meta", metadata.Name, data)

	if err != nil {
		return responseMeta, err
	}

	return ParseShelfResponse(response)
}

// Create metadata property. Will not update existing.
func (this *ShelfLib) CreateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	data, _ := MarshalRequestData(metadata)
	response, _ := this.Request.DoRequest("POST", path, "meta", metadata.Name, data)

	return ParseShelfResponse(response)
}
