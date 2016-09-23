package shelflib_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/quantumew/shelflib"
	"io"
	"log"
	"net/http"
	"path"
)

var validToken = "VALIDTOKEN"
var host = "https://api.shelf.cwscloud.net/"
var logOutput io.Writer
var logger = log.New(logOutput, "", 0)
var shelf = shelflib.New(validToken, *logger)
var testBucket = "test"
var testPath = "test-artifact"

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
		// Get artifact mocked route
		//response := httpmock.NewBytesResponder(200, []byte("simple text"))
		httpmock.RegisterResponder("GET", uriMap["artifact"], func(request *http.Request) (*http.Response, error) {
			token := request.Header["Authorization"][0]

			if token == validToken {
				return httpmock.NewStringResponse(200, "Simple Text File"), nil
			} else {
				return httpmock.NewStringResponse(403, `{"message": "Permission denied", "code": "permission_denied"}`), nil
			}
		})

		httpmock.RegisterResponder("POST", uriMap["artifact"], func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(201, ""), nil
		})

		httpmock.RegisterResponder("PUT", uriMap["meta"], func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
					201,
					`{"version": {"value": "1.5", "immutable": false}, "build": {"value": "10", "immutable": false}}}`),
				nil
		})
	})

	Describe("Integration tests for shelflib", func() {
		Context("GetArtifact", func() {
			It("should successfully retrieve artifact", func() {
				res, err := shelf.GetArtifact(uriMap["artifact"])
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal([]byte("Simple Text File")))
			})

			It("should fail with invalid token", func() {
				shelf.Request.ShelfToken = "INVALID"
				_, err := shelf.GetArtifact(uriMap["artifact"])
				shelfErr := err.(*shelflib.ShelfError)
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

		Context("UpdateMetadataProperty", func() {
			It("should successfully update artifact's metadata property", func() {
			})
		})

		Context("GetMetadata", func() {
			It("should successfully retrieve artifact's metadata", func() {
			})
		})

		Context("GetMetadataProperty", func() {
			It("should successfully retrieve artifact's metadata property", func() {
			})
		})
	})
})
