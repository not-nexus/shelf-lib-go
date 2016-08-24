package shelflib

import (
	"net/http"
    "net/url"
    "path"
)

var SuffixMap = map[string]string{"meta": "_meta", "search":"_search", "artifact":""}

// ShelfRequest encapsulates making requests to shelf-api
type Request struct {
    config Config
}

func (shelfReq Request) Do(verb string, refName string, path string, requestType string) (http.Response, error) {
    requestURI, err := shelfReq.buildUrl(refName, path, requestType)

    if err != nil {
        return nil, err
    }

    client := &http.Client{}
    req, err := http.NewRequest(verb, requestURI, nil)

    if err != nil {
        return nil, err
    }

    req.Header.Add("Authorization", ShelfToken)

    return client.Do(req)
}

func (shelfReq Request) buildUrl(refName string, path string, requestType string) (string, error) {
    url, err := url.Parse(shelfReq.config.ShelfHost)

    if err != nil {
        return nil, err
    }

    suffix := SuffixMap[requestType]
    url.Path := path.Join(refName, shelfReq.config.ShelfPathConst, path, suffix)

    return url.String()
}
