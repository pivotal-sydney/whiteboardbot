package spec

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/model"
)

var _ = Describe("QuietWhiteboard", func() {

	var (
		whiteboard    QuietWhiteboardApp
		store         MockStore
		sydneyStandup Standup
	)

	BeforeEach(func() {
		sydneyStandup = Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"}

		restClient := MockRestClient{}
		restClient.SetStandup(sydneyStandup)

		store = MockStore{}
		whiteboard = NewQuietWhiteboard(&restClient, &store)
	})

	Describe("Receives command", func() {
		Context("?", func() {
			It("should return the usage text", func() {
				expected := Response{Text: USAGE}
				Expect(whiteboard.HandleInput("?")).To(Equal(expected))
			})
		})

		Context("register", func() {
			It("stores the standup in the store", func() {
				expectedStandupJson, _ := json.Marshal(sydneyStandup)
				expectedStandupString := string(expectedStandupJson)

				whiteboard.HandleInput("register 1")

				standupString, standupPresent := store.Get("1")
				Expect(standupPresent).To(Equal(true))
				Expect(standupString).To(Equal(expectedStandupString))
			})

			It("returns a message with the registered standup", func() {
				expected := Response{Text: "Standup Sydney has been registered! You can now start creating Whiteboard entries!"}
				Expect(whiteboard.HandleInput("register 1")).To(Equal(expected))
			})

			Context("when standup does not exist", func() {
				It("returns a message that the standup isn't found", func() {
					expected := Response{Text: "Standup not found!"}
					Expect(whiteboard.HandleInput("register 123")).To(Equal(expected))
				})

				It("does not store anything in the store", func() {
					whiteboard.HandleInput("register 123")
					Expect(len(store.StoreMap)).To(Equal(0))
				})
			})
		})
	})
})
