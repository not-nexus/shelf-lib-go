package shelflib_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
    "github.com/quantumew/shelflib"
)

var _ = Describe("Shelflib", func() {
    Describe("Integration test", func() {
        Context("GetArtifact", func() {
            It("should successfully retrieve artifact", func() {
                _, err := shelflib.GetArtifact("test", "/thing", "TESTTOKEN")
                Expect(err).ShouldNot(HaveOccurred())
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
