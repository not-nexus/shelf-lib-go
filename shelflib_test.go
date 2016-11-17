package shelflib_test

import (
	"github.com/jarcoal/httpmock"
	"github.com/not-nexus/shelf-lib-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tomnomnom/linkheader"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

var validToken = "VALIDTOKEN"
var host = "https://api.shelf.cwscloud.net/"
var logOutput io.Writer
var logger = log.New(logOutput, "", 0)
var shelf = shelflib.New(validToken, logger)
var testBucket = "test"
var testPath = "test-artifact"
var testLink = `</test/artifact/thing>; rel="self"; title="artifact"`
var metadataLink = `</test/artifact/thing/_meta>; rel="related"; title="metadata"`
var testMetadata = map[string]map[string]interface{}{
	"version": map[string]interface{}{"value": "1.5", "immutable": false},
	"build":   map[string]interface{}{"value": "10", "immutable": false},
}

func buildUri(refName string, artifactPath string, requestType string, property string) string {
	suffix := shelflib.SuffixMap[requestType]
	uri := host + path.Join(refName, artifactPath, suffix, property)

	return uri
}

var uriMap = map[string]string{
	"artifact":   buildUri(testBucket, testPath, "artifact", ""),
	"meta":       buildUri(testBucket, testPath, "meta", ""),
	"search":     buildUri(testBucket, testPath, "search", ""),
	"baseSearch": buildUri(testBucket, "", "search", ""),
}

var _ = Describe("Shelflib", func() {
	BeforeEach(func() {
		propResponse := map[string]interface{}{"name": "version", "value": "1.5", "immutable": false}
		permissionsError := map[string]string{"message": "Permission denied", "code": "permission_denied"}
		// Get artifact mocked route
		httpmock.RegisterResponder("GET", uriMap["artifact"], func(request *http.Request) (*http.Response, error) {
			token := request.Header["Authorization"][0]

			if token == validToken {
				return httpmock.NewStringResponse(200, "Simple Text File"), nil
			} else {
				return httpmock.NewJsonResponse(403, permissionsError)
			}
		})

		// Head request for artifact links.
		httpmock.RegisterResponder("HEAD", uriMap["artifact"], func(request *http.Request) (*http.Response, error) {
			response := httpmock.NewStringResponse(204, "")
			response.Header["Link"] = []string{testLink, metadataLink}

			return response, nil
		})

		// Upload artifact.
		httpmock.RegisterResponder("POST", uriMap["artifact"], func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(201, ""), nil
		})

		// Bulk update of metadata.
		httpmock.RegisterResponder("PUT", uriMap["meta"], func(request *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(201, testMetadata)
		})

		// Update metadata property.
		httpmock.RegisterResponder("PUT", uriMap["meta"]+"/version", func(request *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(201, propResponse)
		})

		// Get metadata.
		httpmock.RegisterResponder("GET", uriMap["meta"], func(request *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, testMetadata)
		})

		// Get metadata property.
		httpmock.RegisterResponder("GET", uriMap["meta"]+"/version", func(request *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, propResponse)
		})

		// Create metadata property.
		httpmock.RegisterResponder("POST", uriMap["meta"]+"/stuff", func(request *http.Request) (*http.Response, error) {
			response := map[string]interface{}{"name": "stuff", "value": "monoamine-oxidase-inhibitor", "immutable": true}

			return httpmock.NewJsonResponse(200, response)
		})

		// Search
		httpmock.RegisterResponder("POST", uriMap["search"], func(request *http.Request) (*http.Response, error) {
			response := httpmock.NewStringResponse(204, "")
			response.Header["Link"] = []string{testLink}

			return response, nil
		})
	})

	Describe("Integration tests for shelflib", func() {
		Context("GetArtifact", func() {
			It("should successfully retrieve artifact", func() {
				res, err := shelf.GetArtifact(uriMap["artifact"])
				Expect(err).ShouldNot(HaveOccurred())
				respContents, _ := ioutil.ReadAll(*res)
				Expect(respContents).To(Equal([]byte("Simple Text File")))
			})

			It("should fail with invalid token", func() {
				shelf.Request.ShelfToken = "INVALID"
				_, shelfErr := shelf.GetArtifact(uriMap["artifact"])
				Expect(shelfErr.Message).To(Equal("Permission denied"))
				Expect(shelfErr.Code).To(Equal("permission_denied"))
			})
		})

		Context("CreateArtifact", func() {
			It("should successfully create artifact", func() {
				fileContents := []byte("Simple Text File")
				err := shelf.CreateArtifact(uriMap["artifact"], fileContents)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("UpdateMetadata", func() {
			It("should successfully update artifact's metadata", func() {
				version := &shelflib.MetadataProperty{"version", "1.5", false}
				build := &shelflib.MetadataProperty{"build", "10", false}
				metadata := map[string]*shelflib.MetadataProperty{"version": version, "build": build}
				res, err := shelf.UpdateMetadata(uriMap["artifact"], metadata)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal(metadata))
			})
		})

		Context("ListArtifact", func() {
			It("Should successfuly retrieve artifact link.", func() {
				expectedLinks := linkheader.Parse(testLink)
				links, err := shelf.ListArtifact(uriMap["artifact"])
				Expect(err).ShouldNot(HaveOccurred())
				Expect(*links).To(Equal(expectedLinks))
			})
		})

		Context("UpdateMetadataProperty", func() {
			It("should successfully update artifact's metadata property", func() {
				testProp := &shelflib.MetadataProperty{"version", "1.5", false}
				res, err := shelf.UpdateMetadataProperty(uriMap["artifact"], testProp)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal(testProp))
			})
		})

		Context("GetMetadata", func() {
			It("should successfully retrieve artifact's metadata", func() {
				version := &shelflib.MetadataProperty{"version", "1.5", false}
				build := &shelflib.MetadataProperty{"build", "10", false}
				metadata := map[string]*shelflib.MetadataProperty{"version": version, "build": build}
				res, err := shelf.GetMetadata(uriMap["artifact"])
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal(metadata))
			})
		})

		Context("GetMetadataProperty", func() {
			It("should successfully retrieve artifact's metadata property", func() {
				res, err := shelf.GetMetadataProperty(uriMap["artifact"], "version")
				version := &shelflib.MetadataProperty{"version", "1.5", false}
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal(version))
			})
		})
		Context("CreateMetadataProperty", func() {
			It("successfully creates a metadata property", func() {
				metadata := &shelflib.MetadataProperty{"stuff", "monoamine-oxidase-inhibitor", true}
				res, err := shelf.CreateMetadataProperty(uriMap["artifact"], metadata)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal(metadata))
			})
		})
		Context("Search", func() {
			It("successfully searches", func() {
				expectedLinks := linkheader.Parse(testLink)
				searchCriteria := &shelflib.SearchCriteria{}
				searchCriteria.Search = []string{"artifactName=test-artifact"}
				res, err := shelf.Search(uriMap["artifact"], searchCriteria)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(*res).To(Equal(expectedLinks))
			})
		})
	})
})
