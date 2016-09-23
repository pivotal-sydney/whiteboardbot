package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/model"
)

var _ = Describe("WhiteboardGateway", func() {
	Describe("SaveEntry", func() {
		var (
			restClient MockRestClient
			gateway    WhiteboardGateway
		)

		BeforeEach(func() {
			restClient = MockRestClient{}
			gateway = WhiteboardGateway{RestClient: &restClient}
		})

		It("returns a PostResult with the item ID", func() {
			result, _ := gateway.SaveEntry(&Entry{})

			Expect(result).To(Equal(PostResult{ItemId: "1"}))
		})

		Context("when posting to whiteboard fails", func() {
			It("returns an error with the correct message", func() {
				restClient.SetPostError()

				_, err := gateway.SaveEntry(&Entry{})

				Expect(err.Error()).To(Equal("Problem creating post."))
			})
		})
	})
})
