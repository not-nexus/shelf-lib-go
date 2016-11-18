package shelflib

import (
	"github.com/tomnomnom/linkheader"
	"io"
	"log"
	"os"
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
func (this *ShelfLib) DownloadArtifact(path string) (*io.ReadCloser, *ShelfError) {
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

	return &response.Body, err
}

// Downloads artifact to a file.
func (this *ShelfLib) DownloadArtifactToFile(path string, filePath string) *ShelfError {
	resp, shelfErr := this.DownloadArtifact(path)

	if shelfErr != nil {
		return shelfErr
	}

	outFile, err := os.Create(filePath)

	if err != nil {
		shelfErr = CreateShelfErrorFromError(err)

		return shelfErr
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, *resp)

	if err != nil {
		shelfErr = CreateShelfErrorFromError(err)

		return shelfErr
	}

	return shelfErr
}

// Perform a HEAD request on an artifact endpoint.
// It explicitly REMOVES metadata links.
func (this *ShelfLib) ListArtifact(path string) (*linkheader.Links, *ShelfError) {
	var (
		links linkheader.Links
		err   *ShelfError
	)

	response, err := this.Request.DoRequest("HEAD", path, "artifact", "", nil)

	if err != nil {
		return &links, err
	}

	links, err = ParseLinks(response)

	if err != nil {
		return &links, err
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

	return &links, err
}

// Upload an artifact from Shelf.
func (this *ShelfLib) UploadArtifact(path string, reader io.Reader) *ShelfError {
	response, err := this.Request.Upload(path, reader)

	if err != nil {
		return err
	}

	return CheckResponseStatus(response)
}

// Upload an artifact from a file path.
func (this *ShelfLib) UploadArtifactFromFile(path string, filePath string) *ShelfError {
	reader, err := os.Open(filePath)

	if err != nil {
		shelfErr := CreateShelfErrorFromError(err)

		return shelfErr
	}

	return this.UploadArtifact(path, reader)
}

// Search Shelf using SearchCriteria wrapper struct.
func (this *ShelfLib) Search(path string, searchCriteria *SearchCriteria) (*linkheader.Links, *ShelfError) {
	var links linkheader.Links

	data, err := this.Request.MarshalRequestData(searchCriteria)

	if err != nil {
		return &links, err
	}

	response, err := this.Request.DoRequest("POST", path, "search", "", data)

	if err != nil {
		return &links, err
	}

	links, err = ParseLinks(response)

	return &links, err
}

// Retrieve metadata for an artifact.
func (this *ShelfLib) GetMetadata(path string) (map[string]*MetadataProperty, *ShelfError) {
	var responseMeta map[string]*MetadataProperty

	response, err := this.Request.DoRequest("GET", path, "meta", "", nil)

	if err != nil {
		return responseMeta, err
	}

	return ParseBulkMetadataResponse(response)
}

// Retrieve metadata property for an artifact.
func (this *ShelfLib) GetMetadataProperty(path string, propertyKey string) (*MetadataProperty, *ShelfError) {
	var responseMeta *MetadataProperty

	response, err := this.Request.DoRequest("GET", path, "meta", propertyKey, nil)

	if err != nil {
		return responseMeta, err
	}

	return ParseMetadataResponse(response)
}

// Bulk update of an artifacts metadata.
func (this *ShelfLib) UpdateMetadata(path string, metadata map[string]*MetadataProperty) (map[string]*MetadataProperty, *ShelfError) {
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
func (this *ShelfLib) UpdateMetadataProperty(path string, metadata *MetadataProperty) (*MetadataProperty, *ShelfError) {
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
func (this *ShelfLib) CreateMetadataProperty(path string, metadata *MetadataProperty) (*MetadataProperty, *ShelfError) {
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
