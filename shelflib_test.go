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
		shelfLib := shelflib.New(validToken, logger)

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
				res, err := shelfLib.GetArtifact(uri)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal([]byte("Simple Text File")))
			})

			It("should fail with invalid token", func() {
				uri := host + "test/artifact/thing"
				shelfLib.ShelfToken = "Whatever"
				res, err := shelfLib.GetArtifact(uri)
				Expect(err.Message).To(Equal("Permission denied"))
				Expect(err.Code).To(Equal("permission_denied"))
				Expect(res).To(Equal(nil))
			})
		})

		Context("CreateArtifact", func() {
			It("should successfully create artifact", func() {
				uri := host + "test/artifact/thing"
				fileContents := []byte("Simple Text File")
				res, err := shelfLib.CreateArtifact(uri, fileContents)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("UpdateMetadata", func() {
			It("should successfully update artifact's metadata", func() {
				uri := host + "test/artifact/thing"
				data := map[string]map[string]interface{}{"tag": map[string]interface{}{"value": "val", "immutable": false}}
				res, err := shelflib.UpdateMetadata(uri, data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal(data))
			})
		})

		Context("UpdateMetadataProperty", func() {
			It("should successfully update artifact's metadata property", func() {
				uri := host + "test/artifact/thing"
				data := map[string]interface{}{"value": "val", "immutable": false}
				res, err := shelflib.UpdateMetadata(uri, data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).To(Equal(data))
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
