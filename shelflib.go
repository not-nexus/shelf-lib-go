package shelflib

import (
	"bytes"
	"log"
)

// SearchCriteria is a wrapper for Shelf search criteria.
type SearchCriteria struct {
	Search []string
	Sort   []string
	Limit  int
}

// MetadataProperty is a wrapper for Shelf metadata property.
type MetadataProperty struct {
	Name      string
	Value     string
	Immutable bool
}

type shelfLib struct {
	Logger     log.Logger
	ShelfToken string
}

func New(shelfToken string, logger log.Logger) *shelfLib {
	return &shelfLib{logger, shelfToken}
}

func (this *shelfLib) GetArtifact(path string) ([]byte, error) {
	var artifact []byte
	response, err := DoRequest("GET", this.ShelfToken, path, "artifact", "", nil)

	if err != nil {
		return artifact, err
	}

	resp, err := ParseShelfResponse(response)

	return resp.([]byte), err
}

func (this *shelfLib) ListArtifact(path string) ([]string, error) {
	var links []string
	response, err := DoRequest("HEAD", this.ShelfToken, path, "artifact", "", nil)

	if err != nil {
		return links, err
	}

	resp, err := ParseShelfResponse(response)

	return resp.([]string), err
}

func (this *shelfLib) CreateArtifact(path string, data []byte) error {
	response, err := DoRequest("POST", this.ShelfToken, path, "artifact", "", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	_, err := ParseShelfResponse(response)

	return err
}

func (this *shelfLib) Search(path string, searchCriteria *SearchCriteria) ([]string, error) {
	var links []string
	data, err := MarshalRequestData(searchCriteria)

	if err != nil {
		return links, err
	}

	response, err := DoRequest("POST", this.ShelfToken, path, "search", "", data)

	if err != nil {
		return links, err
	}

	resp, err := ParseShelfResponse(response)

	return resp.([]string), err
}

func (this *shelfLib) GetMetadata(path string) (map[string]MetadataProperty, error) {
	var responseMeta map[string]*MetadataProperty
	response, err := DoRequest("GET", this.ShelfToken, path, "meta", "", nil)

	if err != nil {
		return responseMeta, nil
	}

	resp, err := ParseShelfResponse(response)

	return resp.(map[string]MetadataProperty), err
}

func (this *shelfLib) GetMetadataProperty(path string, property string) (*MetadataProperty, error) {
	var responseMeta MetadataProperty
	response, err := DoRequest("GET", this.ShelfToken, path, "meta", property, nil)

	if err != nil {
		return responseMeta, err
	}

	resp, err := ParseShelfResponse(response)

	return resp.(MetadataProperty), err
}

func (this *shelfLib) UpdateMetadata(path string, metadata map[string]*MetadataProperty) (map[string]*MetadataProperty, error) {
	var responseMeta map[string]*MetadataProperty
	data, err := MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := DoRequest("PUT", this.ShelfToken, path, "meta", "", data)

	if err != nil {
		return responseMeta, err
	}

	return ParseShelfResponse(response)
}

func (this *shelfLib) UpdateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	var responseMeta MetadataProperty

	data, err := MarshalRequestData(metadata)

	if err != nil {
		return responseMeta, err
	}

	response, err := DoRequest("PUT", self.ShelfToken, path, "meta", metadata.Name, data)

	if err != nil {
		return responseMeta, err
	}

	return ParseShelfResponse(response)
}

func (this *shelfLib) CreateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	data, _ := MarshalRequestData(metadata)
	response, _ := DoRequest("POST", self.ShelfToken, path, "meta", metadata.Name, data)

	return ParseShelfResponse(response)
}
