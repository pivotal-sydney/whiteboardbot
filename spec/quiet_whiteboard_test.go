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
		restClient    MockRestClient
	)

	BeforeEach(func() {
		sydneyStandup = Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"}

		user := SlackUser{Username: "aleung", Author: "Andrew Leung"}
		channel := SlackChannel{ChannelId: "C456", ChannelName: "sydney-standup"}
		context = SlackContext{User: user, Channel: channel}

		clock = MockClock{}
		restClient = MockRestClient{}
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
					_, err := whiteboard.ProcessCommand("register 123", context)
					Expect(err.Error()).To(Equal("Standup not found!"))
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

				result, _ := whiteboard.ProcessCommand("faces Nicholas Cage", context)

				Expect(result.Entry).To(Equal(expectedEntry))
			})

			It("stores the new face entry in the entry map", func() {
				title := "Nicholas Cage"
				author := context.User.Author
				expectedEntry := Face{Entry: NewEntry(clock, author, title, sydneyStandup, "New face")}

				whiteboard.ProcessCommand("faces Nicholas Cage", context)

				entry := whiteboard.EntryMap[context.User.Username]
				Expect(entry).To(Equal(expectedEntry))
			})

			It("creates a post", func() {
				whiteboard.ProcessCommand("faces Nicholas Cage", context)
				whiteboard.PostEntry(context)

				expectedRequest := WhiteboardRequest{
					Utf8:   "",
					Method: "",
					Token:  "",
					Item: Item{
						StandupId:   1,
						Title:       "Nicholas Cage",
						Date:        "2015-01-02",
						PostId:      "",
						Public:      "false",
						Kind:        "New face",
						Description: "",
						Author:      "Andrew Leung",
					},
					Commit: "Create New Face",
					Id:     "",
				}

				Expect(restClient.Request).To(Equal(expectedRequest))
				Expect(restClient.PostCalledCount).To(Equal(1))
			})

			Context("when no arguments given", func() {
				It("returns an error message", func() {
					errorMsg := THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"

					expectedEntry := InvalidEntry{Error: errorMsg}

					result, _ := whiteboard.ProcessCommand("faces", context)

					Expect(result.Entry).To(Equal(expectedEntry))

				})

				It("doesn't store anything in the entry map", func() {
					Expect(whiteboard.EntryMap[context.User.Username]).To(BeNil())
					whiteboard.ProcessCommand("faces", context)
					Expect(whiteboard.EntryMap[context.User.Username]).To(BeNil())
				})

				It("doesn't create a post", func() {
					whiteboard.ProcessCommand("faces", context)
					whiteboard.PostEntry(context)

					Expect(restClient.Request).To(Equal(WhiteboardRequest{}))
					Expect(restClient.PostCalledCount).To(BeZero())
				})
			})
		})
	})
})
