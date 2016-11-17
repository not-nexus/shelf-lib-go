package shelflib

import (
	"bytes"
	"github.com/tomnomnom/linkheader"
	"io"
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
	Logger  *log.Logger
	Request *Request
}

// Create a ShelfLib instance.
func New(shelfToken string, logger *log.Logger) *ShelfLib {
	request := &Request{Logger: logger, ShelfToken: shelfToken}

	return &ShelfLib{Logger: logger, Request: request}
}

// Download artifact from Shelf.
func (this *ShelfLib) GetArtifact(path string) (io.ReadCloser, error) {
	response, err := this.Request.DoRequest("GET", path, "artifact", "", nil)

	if err != nil {
		return nil, err
	}

	// Ensures an error response was not returned
	// then it returns the raw body.
	err = CheckResponseStatus(response)

	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

// Perform a HEAD request on an artifact endpoint.
// It explicitly REMOVES metadata links.
func (this *ShelfLib) ListArtifact(path string) (linkheader.Links, error) {
	var links linkheader.Links

	response, err := this.Request.DoRequest("HEAD", path, "artifact", "", nil)

	if err != nil {
		return links, err
	}

	links, err = ParseLinks(response)

	if err != nil {
		return links, err
	}

	for i, link := range links {
		if title, ok := link.Params["title"]; ok {
			if title == "metadata" {
				copy(links[i:], links[i+1:])
				links[len(links)-1] = linkheader.Link{}
				links = links[:len(links)-1]
			}
		}
	}

	return links, nil
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
func (this *ShelfLib) Search(path string, searchCriteria *SearchCriteria) (linkheader.Links, error) {
	var links linkheader.Links

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
		return responseMeta, err
	}

	return ParseBulkMetadataResponse(response)
}

// Retrieve metadata property for an artifact.
func (this *ShelfLib) GetMetadataProperty(path string, propertyKey string) (*MetadataProperty, error) {
	var responseMeta *MetadataProperty

	response, err := this.Request.DoRequest("GET", path, "meta", propertyKey, nil)

	if err != nil {
		return responseMeta, err
	}

	return ParseMetadataResponse(response)
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

	return ParseBulkMetadataResponse(response)
}

// Update metadata property for an artifact.
func (this *ShelfLib) UpdateMetadataProperty(path string, metadata *MetadataProperty) (*MetadataProperty, error) {
	var responseMeta *MetadataProperty

	data, err := this.Request.MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := this.Request.DoRequest("PUT", path, "meta", metadata.Name, data)

	if err != nil {
		return responseMeta, err
	}

	return ParseMetadataResponse(response)
}

// Create metadata property. Will not update existing.
func (this *ShelfLib) CreateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	var responseMeta *MetadataProperty

	data, err := this.Request.MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := this.Request.DoRequest("POST", path, "meta", metadata.Name, data)

	if err != nil {
		return responseMeta, err
	}

	return ParseMetadataResponse(response)
}
