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
	shelfResponse := DoRequest("GET", this.ShelfToken, path, "artifact", "", nil)
}

func (this *shelfLib) ListArtifact(path string) ([]string, error) {
	shelfResponse := DoRequest("HEAD", this.ShelfToken, path, "artifact", "", nil)
}

func (this *shelfLib) CreateArtifact(path string, data []byte) error {
	shelfResponse := DoRequest("POST", this.ShelfToken, path, "artifact", "", bytes.NewBuffer(data))
}

func (this *shelfLib) Search(path string, searchCriteria *SearchCriteria) ([]string, error) {
	data, err := MarshalRequestData(searchCriteria)

	if err != nil {
		return &ShelfResponse{}, err
	}

	shelfResponse := DoRequest("POST", this.ShelfToken, path, "search", "", data)
}

func (this *shelfLib) GetMetadata(path string) (map[string]MetadataProperty, error) {
	shelfResponse := DoRequest("GET", this.ShelfToken, path, "meta", "", nil)
}

func (this *shelfLib) GetMetadataProperty(path string, property string) (*MetadataProperty, error) {
	shelfResponse := DoRequest("GET", this.ShelfToken, path, "meta", property, nil)
}

func (this *shelfLib) UpdateMetadata(path string, metadata map[string]MetadataProperty) (map[string]MetadataProperty, error) {
	data, _ := MarshalRequestData(metadata)

	shelfResponse := DoRequest("PUT", this.ShelfToken, path, "meta", "", data)
}

func (this *shelfLib) UpdateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	data, _ := MarshalRequestData(metadata)
	shelfResponse := oRequest("PUT", self.ShelfToken, path, "meta", metadata.Name, data)
}

func (this *shelfLib) CreateMetadataProperty(path string, metadata MetadataProperty) (*MetadataProperty, error) {
	data, _ := MarshalRequestData(metadata)
	shelfResponse := DoRequest("POST", self.ShelfToken, path, "meta", metadata.Name, data)
}
