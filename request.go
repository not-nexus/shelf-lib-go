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
    ShelfToken string
}

func (shelfReq Request) Do(verb string, refName string, urlPath string, requestType string) (*http.Response, error) {
    response := &http.Response{}
    requestURI, err := shelfReq.buildUrl(refName, urlPath, requestType)

    if err != nil {
        return response, err
    }

    req, err := http.NewRequest(verb, requestURI, nil)

    if err != nil {
        return response, err
    }

    req.Header.Add("Authorization", shelfReq.ShelfToken)
    client := &http.Client{}
    response, err = client.Do(req)

    return response, err
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
