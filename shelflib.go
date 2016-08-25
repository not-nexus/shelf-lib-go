package shelflib

import (
    "bytes"
)

// SearchCriteria is a wrapper for shelf search criteria
type SearchCriteria struct {
    Search []string
    Sort []string
    Limit int
}

var config, _ = LoadConfig()

func GetArtifact(shelfToken string, refName string, path string) (*ShelfResponse, error){
    var request = Request{config, shelfToken}

    return request.Do("GET", refName, path, "artifact", "", nil)
}

func CreateArtifact(shelfToken string, refName string, path string, data []byte) (*ShelfResponse, error) {
    request := Request{config, shelfToken}
    return request.Do("POST", refName, path, "artifact", "", bytes.NewBuffer(data))
}

func Search(shelfToken string, refName string, path string, searchCriteria map[string]interface{}) (*ShelfResponse, error) {
    request := Request{config, shelfToken}
    data, err := request.MarshalRequestData(searchCriteria)

    if err != nil {
        return &ShelfResponse{}, err
    }

    return request.Do("POST", refName, path, "search", "", data)
}

func GetMetadata(shelfToken string, refName string, path string) (*ShelfResponse, error) {
    request := Request{config, shelfToken}
    return request.Do("GET", refName, path, "meta", "", nil)
}

func GetMetadataProperty(shelfToken string, refName string, path string, property string) (*ShelfResponse, error) {
    request := Request{config, shelfToken}
    return request.Do("GET", refName, path, "meta", property, nil)
}

func UpdateMetadata(shelfToken string, refName string, path string, metadata map[string]interface{}) (*ShelfResponse, error) {
    request := Request{config, shelfToken}
    data, _ := request.MarshalRequestData(metadata)

    return request.Do("PUT", refName, path, "meta", "", data)
}

func UpdateMetadataProperty(shelfToken string, refName string, path string, metadata map[string]interface{}, property string) (*ShelfResponse, error) {
    request := Request{config, shelfToken}
    data, _ := request.MarshalRequestData(metadata)
    return request.Do("PUT", refName, path, "meta", property, data)
}

func CreateMetadataProperty(shelfToken string, refName string, path string, metadata map[string]interface{}, property string) (*ShelfResponse, error) {
    request := Request{config, shelfToken}
    data, _ := request.MarshalRequestData(metadata)
    return request.Do("POST", refName, path, "meta", property, data)
}
