package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Rest Client", func() {

	var (
		client MockRestClient
	)

	BeforeEach(func() {
		client = MockRestClient{}
	})

})