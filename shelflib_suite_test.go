package shelflib_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
	"testing"
)

func TestShelflib(t *testing.T) {
	// Activate httpmock for mocking http layer
	// Mock responses are setup in tests themselves.
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Shelflib Suite")
}
