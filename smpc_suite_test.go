package smpc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSmpcGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SmpcGo Suite")
}