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
		context       SlackContext
		clock         MockClock
	)

	BeforeEach(func() {
		sydneyStandup = Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"}

		user := SlackUser{Username: "aleung", Author: "Andrew Leung"}
		channel := SlackChannel{ChannelId: "C456", ChannelName: "sydney-standup"}
		context = SlackContext{User: user, Channel: channel}

		clock = MockClock{}
		restClient := MockRestClient{}
		restClient.SetStandup(sydneyStandup)

		store = MockStore{}
		whiteboard = NewQuietWhiteboard(&restClient, &store, &clock)
	})

	Describe("Receives command", func() {
		Context("?", func() {
			It("should return the usage text", func() {
				expected := CommandResult{Text: USAGE}
				Expect(whiteboard.ProcessCommand("?", context)).To(Equal(expected))
			})
		})

		Context("register", func() {
			It("stores the standup in the store", func() {
				expectedStandupJson, _ := json.Marshal(sydneyStandup)
				expectedStandupString := string(expectedStandupJson)

				whiteboard.ProcessCommand("register 1", context)

				standupString, standupPresent := store.Get("C456")
				Expect(standupPresent).To(Equal(true))
				Expect(standupString).To(Equal(expectedStandupString))
			})

			It("returns a message with the registered standup", func() {
				expected := CommandResult{Text: "Standup Sydney has been registered! You can now start creating Whiteboard entries!"}
				Expect(whiteboard.ProcessCommand("register 1", context)).To(Equal(expected))
			})

			Context("when standup does not exist", func() {
				It("returns a message that the standup isn't found", func() {
					expected := CommandResult{Text: "Standup not found!"}
					Expect(whiteboard.ProcessCommand("register 123", context)).To(Equal(expected))
				})

				It("does not store anything in the store", func() {
					whiteboard.ProcessCommand("register 123", context)
					Expect(len(store.StoreMap)).To(Equal(0))
				})
			})
		})

		Context("faces", func() {
			BeforeEach(func() {
				whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
			})

			It("contains a new face entry in the result", func() {
				title := "Nicholas Cage"
				author := context.User.Author
				expectedEntry := *NewEntry(clock, author, title, sydneyStandup, "New face")

				result := whiteboard.ProcessCommand("faces Nicholas Cage", context)

				Expect(result.Entry).To(Equal(expectedEntry))
			})
		})
	})
})
