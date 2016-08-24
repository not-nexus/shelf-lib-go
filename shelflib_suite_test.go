package shelflib_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShelflib(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shelflib Suite")
}
