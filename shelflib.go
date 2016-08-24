package shelflib

import (
    "fmt"
)

// ShelfResponse is a wrapper for a response from shelf.
type ShelfResponse struct {
    Body interface{}
    Links []string
    StatusCode int
}

// SearchCriteria is a wrapper for shelf search criteria
type SearchCriteria struct {
    Search []string
    Sort []string
    Limit int
}

var request = Request{config: LoadConfig()}

func GetArtifact(refName string, path string) (ShelfResponse, error){
    res, err := request.Do("GET", refName, path, "artifact")
    fmt.Println(res)
}

func CreateArtifact(path string, file []byte) {
}

func Search(path string) {
}

func GetMetadata(path string) {
}

func GetMetadataProperty(key string, path string) {
}

func UpdateMetadata(path string, metadata map[string]interface{}) (map[string]interface{}, error) {
}

func UpdateMetadataProperty(key string, path string, metadata map[string]interface{}) (map[string]interface{}, error) {
}
