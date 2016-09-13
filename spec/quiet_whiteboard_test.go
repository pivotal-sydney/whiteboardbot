package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

var _ = Describe("QuietWhiteboard", func() {

	var (
		whiteboard QuietWhiteboardApp
	)

	BeforeEach(func() {
		restClient := MockRestClient{}
		store := MockStore{}
		whiteboard = NewQuietWhiteboard(&restClient, &store)
	})

	Describe("Receives command", func() {
		Context("?", func() {
			It("should return the usage text", func() {
				expected := Response{Text: USAGE}
				Expect(whiteboard.HandleInput("?")).To(Equal(expected))
			})
		})
	})
})
