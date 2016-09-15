package shelflib_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/quantumew/shelflib"
	"net/http"
	"os"
)

var validToken = "VALIDTOKEN"
var host = "https://api.shelf.cwscloud.net/"

var _ = Describe("Shelflib", func() {
	BeforeEach(func() {
		var logOutput []byte
		logger := log.New(logOutput, "", 0)
		validShelf := shelflib.New(validToken, logger)
		invalidTokenShelf := shelflib.New("WHATEVER", logger)

		// Get artifact mocked route
		//response := httpmock.NewBytesResponder(200, []byte("simple text"))
		httpmock.RegisterResponder("GET", host+"test/artifact/thing", func(request *http.Request) (*http.Response, error) {
			token := request.Header["Authorization"][0]

			if token == validToken {
				return httpmock.NewStringResponse(200, "Simple Text File"), nil
			} else {
				return httpmock.NewStringResponse(403, `{"message": "Permission denied", "code": "permission_denied"}`), nil
			}
		})

		httpmock.RegisterResponder("POST", host+"test/artifact/thing", func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(201, ""), nil
		})

		httpmock.RegisterResponder("PUT", host+"test/artifact/thing/_meta", func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(201, `{"tag": {"value": "val", "immutable": false}}}`), nil
		})
	})

	Describe("Integration test", func() {
		Context("GetArtifact", func() {
			It("should successfully retrieve artifact", func() {
				uri := host + "test/artifact/thing"
				res, err := validShelf.GetArtifact(uri)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(200))
			})

			It("should fail with invalid token", func() {
				uri := host + "test/artifact/thing"
				res, err := invalidTokenShelf.GetArtifact(uri)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(403))
			})
		})

		Context("CreateArtifact", func() {
			It("should successfully create artifact", func() {
				fileContents := []byte("Simple Text File")
				res, err := shelflib.CreateArtifact("TOKEN", "test", "/thing", fileContents)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(201))
			})
		})

		Context("UpdateMetadata", func() {
			It("should successfully update artifact's metadata", func() {
				data := map[string]map[string]interface{}{"tag": map[string]interface{}{"value": "val", "immutable": false}}
				res, err := shelflib.UpdateMetadata("TOKEN", "test", "/thing", data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(201))
			})
		})

		Context("UpdateMetadataProperty", func() {
			It("should successfully update artifact's metadata property", func() {
				data := map[string]interface{}{"value": "val", "immutable": false}
				res, err := shelflib.UpdateMetadata("TOKEN", "test", "/thing", data, "tag")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(201))
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
