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
				expected := CommandResult{Entry: TextEntry{Text: USAGE}}
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
				expected := CommandResult{Entry: TextEntry{Text: "Standup Sydney has been registered! You can now start creating Whiteboard entries!"}}
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

		Context("body", func() {
			BeforeEach(func() {
				whiteboard.ProcessCommand("register 1", context)
			})

			It("adds a body to the entry type", func() {
				whiteboard.ProcessCommand("interestings Nicholas Cage did a remake of The Wicker Man!", context)
				entryType := whiteboard.EntryMap[context.User.Username]
				Expect(entryType.GetEntry().Body).To(BeEmpty())

				whiteboard.ProcessCommand("body And the movie was terrible!", context)
				entryType = whiteboard.EntryMap[context.User.Username]

				Expect(entryType.GetEntry().Body).To(Equal("And the movie was terrible!"))
			})

			It("creates a post", func() {
				whiteboard.ProcessCommand("interestings Nicholas Cage did a remake of The Wicker Man!", context)
				whiteboard.ProcessCommand("body And the movie was terrible!", context)

				expectedRequest := WhiteboardRequest{
					Utf8:   "",
					Method: "",
					Token:  "",
					Item: Item{
						StandupId:   1,
						Title:       "Nicholas Cage did a remake of The Wicker Man!",
						Date:        "2015-01-02",
						PostId:      "",
						Public:      "false",
						Kind:        "Interesting",
						Description: "And the movie was terrible!",
						Author:      "Andrew Leung",
					},
					Commit: "Create Item",
					Id:     "",
				}

				Expect(restClient.Request).To(Equal(expectedRequest))
				Expect(restClient.PostCalledCount).To(Equal(2))
			})

			Context("when there is no entry", func() {
				It("returns an error", func() {
					errorMsg := ":-1:\nHey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!"
					expectedEntry := InvalidEntry{Error: errorMsg}

					result, _ := whiteboard.ProcessCommand("body And the movie was terrible!", context)

					Expect(result).To(Equal(CommandResult{Entry: expectedEntry}))
				})
			})

			Context("when entry type is a New Face", func() {
				FIt("returns an error", func() {
					whiteboard.ProcessCommand("faces Nicholas Cage", context)
					errorMsg := ":-1:\nHey, new faces should not have a body!"
					expectedEntry := InvalidEntry{Error: errorMsg}

					result, _ := whiteboard.ProcessCommand("body And John Travolta!", context)

					Expect(result).To(Equal(CommandResult{Entry: expectedEntry}))
				})
			})

			Context("when given no arguments", func() {
				It("returns an error", func() {
					whiteboard.ProcessCommand("interestings Nicholas Cage did a remake of The Wicker Man!", context)
					errorMsg := THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb b My title`\nNeed help? Try `wb ?`"
					expectedEntry := InvalidEntry{Error: errorMsg}

					result, _ := whiteboard.ProcessCommand("body", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				})
			})

		})

		Describe("creating entries", func() {
			var (
				title                 string
				author                string
				command               string
				commitString          string
				expectedEntry         Entry
				expectedEntryType     EntryType
				expectedEntryItemKind string
			)

			AssertContainsEntryInResult := func() func() {
				return func() {
					result, _ := whiteboard.ProcessCommand(command+" "+title, context)

					Expect(result.Entry).To(Equal(expectedEntry))
				}
			}

			AssertEntryStoredInEntryMap := func() func() {
				return func() {
					whiteboard.ProcessCommand(command+" "+title, context)

					entry := whiteboard.EntryMap[context.User.Username]
					Expect(entry).To(Equal(expectedEntryType))
				}
			}

			AssertPostCreated := func() func() {
				return func() {
					whiteboard.ProcessCommand(command+" "+title, context)

					expectedRequest := WhiteboardRequest{
						Utf8:   "",
						Method: "",
						Token:  "",
						Item: Item{
							StandupId:   1,
							Title:       title,
							Date:        "2015-01-02",
							PostId:      "",
							Public:      "false",
							Kind:        expectedEntry.ItemKind,
							Description: "",
							Author:      author,
						},
						Commit: commitString,
						Id:     "",
					}

					Expect(restClient.Request).To(Equal(expectedRequest))
					Expect(restClient.PostCalledCount).To(Equal(1))
				}
			}

			AssertErrorMessageWhenCreatingPostFails := func() func() {
				return func() {
					restClient.SetPostError()
					expectedEntry := InvalidEntry{Error: "Problem creating post."}

					result, _ := whiteboard.ProcessCommand(command+" "+title, context)

					Expect(result.Entry).To(Equal(expectedEntry))
				}
			}

			AssertNoArgumentErrorMessage := func() func() {
				return func() {
					errorMsg := THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"

					expectedEntry := InvalidEntry{Error: errorMsg}

					result, _ := whiteboard.ProcessCommand(command, context)

					Expect(result.Entry).To(Equal(expectedEntry))
				}
			}

			AssertNoArgumentWontStoreInEntryMap := func() func() {
				return func() {
					Expect(whiteboard.EntryMap[context.User.Username]).To(BeNil())
					whiteboard.ProcessCommand(command, context)
					Expect(whiteboard.EntryMap[context.User.Username]).To(BeNil())
				}
			}

			AssertNoArgumentWontCreatePost := func() func() {
				return func() {
					whiteboard.ProcessCommand(command, context)

					Expect(restClient.Request).To(Equal(WhiteboardRequest{}))
					Expect(restClient.PostCalledCount).To(BeZero())
				}
			}

			Context("faces", func() {
				BeforeEach(func() {
					command = "faces"
					commitString = "Create New Face"
					expectedEntryItemKind = "New face"
					title = "Nicholas Cage"
					author = context.User.Author
					expectedEntry = *NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					expectedEntryType = Face{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("creates a post", AssertPostCreated())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenCreatingPostFails())
				})

				Context("when no arguments given", func() {
					It("returns an error message", AssertNoArgumentErrorMessage())

					It("doesn't store anything in the entry map", AssertNoArgumentWontStoreInEntryMap())

					It("doesn't create a post", AssertNoArgumentWontCreatePost())
				})
			})

			Context("helps", func() {

				BeforeEach(func() {
					command = "helps"
					commitString = "Create Item"
					expectedEntryItemKind = "Help"
					title = "Good wicker furniture shop recommendations?"
					author = context.User.Author
					expectedEntry = *NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					expectedEntryType = Help{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("creates a post", AssertPostCreated())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenCreatingPostFails())
				})

				Context("when no arguments given", func() {
					It("returns an error message", AssertNoArgumentErrorMessage())

					It("doesn't store anything in the entry map", AssertNoArgumentWontStoreInEntryMap())

					It("doesn't create a post", AssertNoArgumentWontCreatePost())
				})
			})

			Context("interestings", func() {

				BeforeEach(func() {
					command = "interestings"
					commitString = "Create Item"
					expectedEntryItemKind = "Interesting"
					title = "Nicholas Cage did a remake of The Wicker Man!"
					author = context.User.Author
					expectedEntry = *NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					expectedEntryType = Interesting{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("creates a post", AssertPostCreated())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenCreatingPostFails())
				})

				Context("when no arguments given", func() {
					It("returns an error message", AssertNoArgumentErrorMessage())

					It("doesn't store anything in the entry map", AssertNoArgumentWontStoreInEntryMap())

					It("doesn't create a post", AssertNoArgumentWontCreatePost())
				})
			})

			Context("events", func() {

				BeforeEach(func() {
					command = "events"
					commitString = "Create Item"
					expectedEntryItemKind = "Event"
					title = "Movie Screening for The Wicker Man!"
					author = context.User.Author
					expectedEntry = *NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					expectedEntryType = Event{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("creates a post", AssertPostCreated())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenCreatingPostFails())
				})

				Context("when no arguments given", func() {
					It("returns an error message", AssertNoArgumentErrorMessage())

					It("doesn't store anything in the entry map", AssertNoArgumentWontStoreInEntryMap())

					It("doesn't create a post", AssertNoArgumentWontCreatePost())
				})
			})
		})
	})
})
