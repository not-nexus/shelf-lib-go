package shelflib

import (
	"net/http"
    "net/url"
    "path"
    "io"
    "encoding/json"
    "bytes"
)

var SuffixMap = map[string]string{"meta": "_meta", "search":"_search", "artifact":""}

// ShelfResponse is a wrapper for a response from shelf.
type ShelfResponse struct {
    Body interface{}
    Links []string
    StatusCode int
}

// ShelfRequest encapsulates making requests to shelf-api
type Request struct {
    config Config
    ShelfToken string
}

func (shelfReq Request) MarshalRequestData(data interface{}) (io.Reader, error) {
    jsonData, err := json.Marshal(data)

    if err != nil {
        return nil, err
    }

    return bytes.NewBuffer(jsonData), nil
}

func (shelfReq Request) Do(verb string, refName string, urlPath string, requestType string, property string, data io.Reader) (*ShelfResponse, error) {
    response := &ShelfResponse{}
    requestURI, err := shelfReq.buildUrl(refName, urlPath, requestType)

    if err != nil {
        return response, err
    }

    req, err := http.NewRequest(verb, requestURI, data)

    if err != nil {
        return response, err
    }

    req.Header.Add("Authorization", shelfReq.ShelfToken)
    client := &http.Client{}
    rawResponse, err := client.Do(req)

    if err != nil {
        return response, err
    }

    return shelfReq.parseRequest(rawResponse)
}

func (shelfReq Request) parseRequest(response *http.Response) (*ShelfResponse, error) {
    shelfResponse := &ShelfResponse{
        StatusCode: response.StatusCode,
        Links: response.Header["Link"],
        Body: response.Body,
    }

    return shelfResponse, nil
}

func (shelfReq Request) buildUrl(refName string, urlPath string, requestType string) (string, error) {
    url, err := url.Parse(shelfReq.config.ShelfHost)

    if err != nil {
        return "", err
    }

    suffix := SuffixMap[requestType]
    url.Path = path.Join(refName, shelfReq.config.ShelfPathConst, urlPath, suffix)

    return url.String(), nil
}
