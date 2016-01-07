package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWhiteboardbot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Whiteboardbot Suite")
}