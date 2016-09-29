package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/model"
)

var _ = Describe("WhiteboardGateway", func() {
	var (
		restClient MockRestClient
		gateway    WhiteboardGateway
	)

	BeforeEach(func() {
		restClient = MockRestClient{}
		gateway = WhiteboardGateway{RestClient: &restClient}
	})

	Describe("FindStandup", func() {
		It("returns the standup", func() {
			expectedStandup := Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"}
			restClient.SetStandup(expectedStandup)

			standup, _ := gateway.FindStandup("1")

			Expect(standup).To(Equal(expectedStandup))
		})

		Context("when the standup is not found", func() {
			It("returns an error message", func() {
				_, err := gateway.FindStandup("abc123")
				Expect(err.Error()).To(Equal("Standup not found!"))
			})
		})
	})

	Describe("SaveEntry", func() {
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

	Describe("GetStandupItems", func() {
		var (
			standup              Standup
			expectedStandupItems StandupItems
		)
		BeforeEach(func() {
			standup = Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"}
			restClient.SetStandup(standup)
			expectedStandupItems = StandupItems{Interestings: []Entry{
				{Title: "I1", Author: "Alice", Date: "2015-01-02"},
			}}
			restClient.StandupItems = expectedStandupItems
		})

		It("invokes the rest client's get standup items method", func() {
			Expect(gateway.GetStandupItems("1")).To(Equal(expectedStandupItems))
		})

		Context("when the standup id is not valid", func() {
			It("returns an error message", func() {
				_, err := gateway.GetStandupItems("foo")

				Expect(err.Error()).To(Equal(MISSING_STANDUP))
			})
		})

		Context("when fetching standup items fails", func() {
			It("returns an error message", func() {
				restClient.SetGetStandupItemsError()
				_, err := gateway.GetStandupItems("1")
				Expect(err.Error()).To(Equal("Failed fetching standup items"))
			})
		})
	})
})
