package shelflib_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
    "github.com/quantumew/shelflib"
    "github.com/jarcoal/httpmock"
    "net/http"
)

var ValidToken = "VALIDTOKEN"

var _ = Describe("Shelflib", func() {
    BeforeEach(func() {
        // Get artifact mocked route
        //response := httpmock.NewBytesResponder(200, []byte("simple text"))
        httpmock.RegisterResponder("GET", "https://api.shelf.cwscloud.net/test/artifact/thing", func(request *http.Request) (*http.Response, error) {
            token := request.Header["Authorization"][0]

            if token == ValidToken {
                return httpmock.NewStringResponse(200, ""), nil
            } else {
                return httpmock.NewStringResponse(403, `{"message": "Permission denied", "code": "permission_denied"}`), nil
            }
        })
    })

    Describe("Integration test", func() {
        Context("GetArtifact", func() {
            It("should successfully retrieve artifact", func() {
                res, err := shelflib.GetArtifact(ValidToken, "test", "/thing")
                Expect(err).ShouldNot(HaveOccurred())
                Expect(res.StatusCode).To(Equal(200))
            })

            It("should fail with invalid token", func() {
                res, err := shelflib.GetArtifact("BLAH", "test", "/thing")
                Expect(err).ShouldNot(HaveOccurred())
                Expect(res.StatusCode).To(Equal(403))
            })
        })

        Context("CreateArtifact", func() {
            It("should successfully create artifact", func() {
            })
        })

        Context("UpdateMetadata", func() {
            It("should successfully update artifact's metadata", func() {
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
