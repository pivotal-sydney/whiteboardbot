package spec

import (
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"strconv"
)

type MockWhiteboardGateway struct {
	StandupMap      map[int]Standup
	SaveEntryCalled bool
	EntrySaved      EntryType
	failSaveEntry   bool
}

func (gateway *MockWhiteboardGateway) FindStandup(standupId string) (standup Standup, err error) {
	var ok bool
	id, _ := strconv.Atoi(standupId)

	standup, ok = gateway.StandupMap[id]

	if !ok {
		err = errors.New("Standup not found!")
	}

	return
}

func (gateway *MockWhiteboardGateway) SaveEntry(entryType EntryType) (PostResult, error) {
	if gateway.failSaveEntry {
		return PostResult{}, errors.New("Problem creating post.")
	}
	gateway.SaveEntryCalled = true
	gateway.EntrySaved = entryType

	return PostResult{ItemId: "1"}, nil
}

func (gateway *MockWhiteboardGateway) SetSaveEntryError() {
	gateway.failSaveEntry = true
}

func (gateway *MockWhiteboardGateway) SetStandup(standup Standup) {
	if gateway.StandupMap == nil {
		gateway.StandupMap = make(map[int]Standup)
	}
	gateway.StandupMap[standup.Id] = standup
}

var _ = Describe("QuietWhiteboard", func() {

	var (
		whiteboard    QuietWhiteboardApp
		store         MockStore
		sydneyStandup Standup
		context       SlackContext
		clock         MockClock
		gateway       MockWhiteboardGateway
	)

	BeforeEach(func() {
		sydneyStandup = Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"}

		user := SlackUser{Username: "aleung", Author: "Andrew Leung"}
		channel := SlackChannel{ChannelId: "C456", ChannelName: "sydney-standup"}
		context = SlackContext{User: user, Channel: channel}

		clock = MockClock{}

		store = MockStore{}

		gateway = MockWhiteboardGateway{}
		gateway.SetStandup(sydneyStandup)

		whiteboard = NewQuietWhiteboard(&MockRestClient{}, &gateway, &store, &clock)
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
					expectedEntry := InvalidEntry{Error: "Standup not found!"}

					result := whiteboard.ProcessCommand("register 123", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				})

				It("does not store anything in the store", func() {
					whiteboard.ProcessCommand("register 123", context)
					Expect(len(store.StoreMap)).To(Equal(0))
				})
			})
		})

		Context("body", func() {
			var expectedEntry EntryType

			BeforeEach(func() {
				entry := Entry{
					Date:      "2015-01-02",
					Title:     "Nicholas Cage did a remake of The Wicker Man!",
					Body:      "",
					Author:    "Andrew Leung",
					Id:        "1",
					StandupId: 1,
					ItemKind:  "Interesting",
				}
				expectedEntry = Interesting{Entry: &entry}
				whiteboard.EntryMap[context.User.Username] = expectedEntry
			})

			It("adds a body to the entry type", func() {
				entryType := whiteboard.EntryMap[context.User.Username]
				Expect(entryType.GetEntry().Body).To(BeEmpty())

				whiteboard.ProcessCommand("body And the movie was terrible!", context)
				entryType = whiteboard.EntryMap[context.User.Username]

				Expect(entryType.GetEntry().Body).To(Equal("And the movie was terrible!"))
			})

			It("creates a post", func() {
				whiteboard.ProcessCommand("body And the movie was terrible!", context)

				Expect(gateway.SaveEntryCalled).To(BeTrue())
				Expect(gateway.EntrySaved).To(Equal(expectedEntry))
			})

			Context("when there is no entry", func() {
				It("returns an error", func() {
					delete(whiteboard.EntryMap, context.User.Username)

					expectedEntry := InvalidEntry{Error: MISSING_ENTRY}

					result := whiteboard.ProcessCommand("body And the movie was terrible!", context)

					Expect(result).To(Equal(CommandResult{Entry: expectedEntry}))
				})
			})

			Context("when entry type is a New Face", func() {
				It("returns an error", func() {
					whiteboard.ProcessCommand("faces Nicholas Cage", context)
					errorMsg := ":-1:\nHey, new faces should not have a body!"
					expectedEntry := InvalidEntry{Error: errorMsg}

					result := whiteboard.ProcessCommand("body And John Travolta!", context)

					Expect(result).To(Equal(CommandResult{Entry: expectedEntry}))
				})
			})

			Context("when given no arguments", func() {
				It("returns an error", func() {
					whiteboard.ProcessCommand("interestings Nicholas Cage did a remake of The Wicker Man!", context)
					expectedEntry := InvalidEntry{Error: MISSING_INPUT}

					result := whiteboard.ProcessCommand("body", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				})
			})

			Context("when saving the entry fails", func() {
				It("returns an error message", func() {
					gateway.SetSaveEntryError()
					expectedEntry := InvalidEntry{Error: "Problem creating post."}

					result := whiteboard.ProcessCommand("body And the movie was terrible!", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				})
			})
		})

		Context("date", func() {
			var expectedEntry EntryType

			BeforeEach(func() {
				entry := Entry{
					Date:      "2015-01-02",
					Title:     "Nicholas Cage did a remake of The Wicker Man!",
					Body:      "",
					Author:    "Andrew Leung",
					Id:        "1",
					StandupId: 1,
					ItemKind:  "Interesting",
				}
				expectedEntry = Interesting{Entry: &entry}
				whiteboard.EntryMap[context.User.Username] = expectedEntry
			})

			It("updates the date of an entry", func() {
				entryType := whiteboard.EntryMap[context.User.Username]
				Expect(entryType.GetDateString()).To(Equal("02 Jan 2015"))

				whiteboard.ProcessCommand("date 3000-05-13", context)
				entryType = whiteboard.EntryMap[context.User.Username]

				Expect(entryType.GetDateString()).To(Equal("13 May 3000"))
			})

			It("updates the entry", func() {
				whiteboard.ProcessCommand("date 3000-05-13", context)

				Expect(gateway.SaveEntryCalled).To(BeTrue())
				Expect(gateway.EntrySaved.GetDateString()).To(Equal("13 May 3000"))
			})

			Context("when there is no entry", func() {
				It("returns an error", func() {
					delete(whiteboard.EntryMap, context.User.Username)

					expectedEntry := InvalidEntry{Error: MISSING_ENTRY}

					result := whiteboard.ProcessCommand("date 3000-05-13", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				})
			})

			Context("when given no arguments", func() {
				It("returns an error", func() {
					expectedEntry := InvalidEntry{Error: MISSING_INPUT}

					result := whiteboard.ProcessCommand("date", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				})
			})

			Context("when the date format is wrong", func() {
				It("returns an error", func() {
					errorMsg := THUMBS_DOWN + "Date not set, use YYYY-MM-DD as date format\n"
					expectedEntry := InvalidEntry{Error: errorMsg}
					whiteboard.ProcessCommand("interestings Nicholas Cage did a remake of The Wicker Man!", context)

					result := whiteboard.ProcessCommand("date LOLWUT", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				})
			})

			Context("when saving the entry fails", func() {
				It("returns an error message", func() {
					gateway.SetSaveEntryError()
					expectedEntry := InvalidEntry{Error: "Problem creating post."}

					result := whiteboard.ProcessCommand("date 3000-05-13", context)

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
					result := whiteboard.ProcessCommand(command+" "+title, context)

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

			AssertEntryHasAnId := func() func() {
				return func() {
					whiteboard.ProcessCommand(command+" "+title, context)

					entry := whiteboard.EntryMap[context.User.Username]
					Expect(entry.GetEntry().Id).To(Equal("1"))
				}
			}

			AssertEntrySaved := func() func() {
				return func() {
					whiteboard.ProcessCommand(command+" "+title, context)

					Expect(gateway.SaveEntryCalled).To(BeTrue())
					Expect(gateway.EntrySaved).To(Equal(expectedEntryType))
				}
			}

			AssertErrorMessageWhenEntrySaveFails := func() func() {
				return func() {
					gateway.SetSaveEntryError()
					expectedEntry := InvalidEntry{Error: "Problem creating post."}

					result := whiteboard.ProcessCommand(command+" "+title, context)

					Expect(result.Entry).To(Equal(expectedEntry))
				}
			}

			AssertNoArgumentErrorMessage := func() func() {
				return func() {
					expectedEntry := InvalidEntry{Error: MISSING_INPUT}

					result := whiteboard.ProcessCommand(command, context)

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

					Expect(gateway.SaveEntryCalled).To(BeFalse())
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
					expectedEntry.Id = "1"
					expectedEntryType = Face{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("assigns an Id to the entry", AssertEntryHasAnId())

				It("creates a post", AssertEntrySaved())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenEntrySaveFails())
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
					expectedEntry.Id = "1"
					expectedEntryType = Help{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("assigns an Id to the entry", AssertEntryHasAnId())

				It("creates a post", AssertEntrySaved())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenEntrySaveFails())
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
					expectedEntry.Id = "1"
					expectedEntryType = Interesting{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("assigns an Id to the entry", AssertEntryHasAnId())

				It("creates a post", AssertEntrySaved())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenEntrySaveFails())
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
					expectedEntry.Id = "1"
					expectedEntryType = Event{Entry: &expectedEntry}

					whiteboard.Store.SetStandup(context.Channel.ChannelId, sydneyStandup)
				})

				It("contains a help entry in the result", AssertContainsEntryInResult())

				It("stores the help entry in the entry map", AssertEntryStoredInEntryMap())

				It("assigns an Id to the entry", AssertEntryHasAnId())

				It("creates a post", AssertEntrySaved())

				Context("when posting to Whiteboard fails", func() {
					It("returns the proper error message", AssertErrorMessageWhenEntrySaveFails())
				})

				Context("when no arguments given", func() {
					It("returns an error message", AssertNoArgumentErrorMessage())

					It("doesn't store anything in the entry map", AssertNoArgumentWontStoreInEntryMap())

					It("doesn't create a post", AssertNoArgumentWontCreatePost())
				})
			})
		})

		Context("updating entries", func() {
			var (
				originalEntry   Entry
				originalEntryId string
				command         string
				newValue        string
			)

			AssertTitleUpdated := func() func() {
				return func() {
					whiteboard.ProcessCommand(command+" "+newValue, context)

					entry := whiteboard.EntryMap[context.User.Username].GetEntry()
					entryId := entry.Id
					title := entry.Title

					Expect(title).To(Equal(newValue))
					Expect(entryId).To(Equal(originalEntryId))
				}
			}

			AssertNoTitleErrorMessage := func() func() {
				return func() {
					expectedEntry := InvalidEntry{Error: MISSING_INPUT}

					result := whiteboard.ProcessCommand(command+"     ", context)

					Expect(result.Entry).To(Equal(expectedEntry))
				}
			}

			AssertNoEntryErrorMessage := func() func() {
				return func() {
					expectedEntry := InvalidEntry{Error: MISSING_ENTRY}

					delete(whiteboard.EntryMap, context.User.Username)

					result := whiteboard.ProcessCommand(command+" "+newValue, context)

					Expect(result.Entry).To(Equal(expectedEntry))
				}
			}

			AssertSaveEntryFailureErrorMessage := func() func() {
				return func() {
					gateway.SetSaveEntryError()
					expectedEntry := InvalidEntry{Error: "Problem creating post."}

					result := whiteboard.ProcessCommand(command+" "+newValue, context)

					Expect(result.Entry).To(Equal(expectedEntry))
				}
			}

			Context("name", func() {
				BeforeEach(func() {
					command = "name"
					newValue = "Olivia Newton John"
					originalEntryId = "abc123"
					originalEntryType := NewFace(clock, context.User.Author, "Oliver Newton John", sydneyStandup)
					originalEntryType.GetEntry().Id = originalEntryId
					whiteboard.EntryMap[context.User.Username] = originalEntryType
					originalEntry = *whiteboard.EntryMap[context.User.Username].GetEntry()
				})

				It("updates the name on a new face", AssertTitleUpdated())

				Context("when the new name is the empty string", func() {
					It("returns an error message", AssertNoTitleErrorMessage())
				})

				Context("no entry in store", func() {
					It("returns an error message", AssertNoEntryErrorMessage())
				})
				Context("when saving the entry fails", func() {
					It("returns an error message", AssertSaveEntryFailureErrorMessage())
				})
			})

			Context("title", func() {
				BeforeEach(func() {
					command = "title"
					newValue = "Saturday Night Live"
					originalEntryId = "abc123"
					originalEntryType := NewEvent(clock, context.User.Author, "Saturday Night Fever", sydneyStandup)
					originalEntryType.GetEntry().Id = originalEntryId
					whiteboard.EntryMap[context.User.Username] = originalEntryType
					originalEntry = *whiteboard.EntryMap[context.User.Username].GetEntry()
				})

				It("updates the title on an entry", AssertTitleUpdated())

				Context("when the new title is the empty string", func() {
					It("returns an error message", AssertNoTitleErrorMessage())
				})

				Context("no entry in store", func() {
					It("returns an error message", AssertNoEntryErrorMessage())
				})
			})
		})
	})
})
