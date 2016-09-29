package spec

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/model"
)

var _ = Describe("EntryCommandResult", func() {
	Describe("String", func() {
		It("prints the status, help text, title and entry string", func() {
			entry := NewEntry(&MockClock{}, "Ernest Hemingway", "Peace and War", Standup{}, "Book")
			entryCommandResult := EntryCommandResult{
				Title:    "Good Book",
				Status:   "SUCCESS!",
				HelpText: "Here's a little help",
				Entry:    entry,
			}

			expectedString := "SUCCESS!\nGood Book\nHere's a little help\n" + entry.String()

			Expect(entryCommandResult.String()).To(Equal(expectedString))
		})

		Context("when the help text is empty", func() {
			It("skips the help text", func() {
				entry := NewEntry(&MockClock{}, "Ernest Hemingway", "Peace and War", Standup{}, "Book")
				entryCommandResult := EntryCommandResult{
					Title:    "Good Book",
					Status:   "SUCCESS!",
					HelpText: "",
					Entry:    entry,
				}

				expectedString := "SUCCESS!\nGood Book\n" + entry.String()

				Expect(entryCommandResult.String()).To(Equal(expectedString))
			})
		})
	})
})

var _ = Describe("MessageCommandResult", func() {
	Describe("String", func() {
		It("prints the status and text", func() {
			messageCommandResult := MessageCommandResult{
				Status: "SUCCESS!",
				Text:   "Some text",
			}

			Expect(messageCommandResult.String()).To(Equal("SUCCESS!\nSome text"))
		})

		Context("when status is empty", func() {
			It("skips the status text", func() {
				messageCommandResult := MessageCommandResult{
					Status: "",
					Text:   "Some text",
				}

				Expect(messageCommandResult.String()).To(Equal("Some text"))
			})
		})
	})
})

var _ = Describe("Whiteboard", func() {

	var (
		whiteboard    WhiteboardApp
		store         MockStore
		sydneyStandup Standup
		context       SlackContext
		clock         MockClock
		gateway       MockWhiteboardGateway
	)

	BeforeEach(func() {
		sydneyStandup = Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"}

		user := SlackUser{Username: "aleung", Author: "Andrew Leung"}
		channel := SlackChannel{Id: "C456", Name: "sydney-standup"}
		context = SlackContext{User: user, Channel: channel}

		clock = MockClock{}

		store = MockStore{}
		store.SetStandup(context.Channel.Id, sydneyStandup)

		gateway = MockWhiteboardGateway{}
		gateway.SetStandup(sydneyStandup)

		whiteboard = NewWhiteboard(&gateway, &store, &clock)
	})

	Describe("Receives command", func() {
		Describe("?", func() {
			It("should return the usage text", func() {
				expected := MessageCommandResult{Text: USAGE}
				Expect(whiteboard.ProcessCommand("?", context)).To(Equal(expected))
			})
		})

		Describe("register", func() {
			BeforeEach(func() {
				store.StoreMap = make(map[string]string)
			})

			It("stores the standup in the store", func() {
				expectedStandupJson, _ := json.Marshal(sydneyStandup)
				expectedStandupString := string(expectedStandupJson)

				whiteboard.ProcessCommand("register 1", context)

				standupString, standupPresent := store.Get("C456")
				Expect(standupPresent).To(Equal(true))
				Expect(standupString).To(Equal(expectedStandupString))
			})

			It("returns a message with the registered standup", func() {
				expected := MessageCommandResult{
					Text:   "Standup Sydney has been registered! You can now start creating Whiteboard entries!",
					Status: THUMBS_UP,
				}
				Expect(whiteboard.ProcessCommand("register 1", context)).To(Equal(expected))
			})

			Context("when standup does not exist", func() {
				It("returns a message that the standup isn't found", func() {
					expectedResult := MessageCommandResult{
						Text:   "Standup not found!",
						Status: THUMBS_DOWN,
					}

					result := whiteboard.ProcessCommand("register 123", context)

					Expect(result).To(Equal(expectedResult))
				})

				It("does not store anything in the store", func() {
					whiteboard.ProcessCommand("register 123", context)
					Expect(len(store.StoreMap)).To(Equal(0))
				})
			})
		})

		Describe("creating entries", func() {
			var (
				title                 string
				author                string
				command               string
				commitString          string
				expectedEntryType     EntryType
				expectedEntryItemKind string
				expectedResult        EntryCommandResult
			)

			AssertContainsEntryInResult := func() func() {
				return func() {
					result := whiteboard.ProcessCommand(command+" "+title, context)

					Expect(result).To(Equal(expectedResult))
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
					expectedResult := MessageCommandResult{Text: "Problem creating post."}

					result := whiteboard.ProcessCommand(command+" "+title, context)

					Expect(result).To(Equal(expectedResult))
				}
			}

			AssertNoArgumentErrorMessage := func() func() {
				return func() {
					expectedResult := MessageCommandResult{Text: MISSING_INPUT, Status: THUMBS_DOWN}

					result := whiteboard.ProcessCommand(command, context)

					Expect(result).To(Equal(expectedResult))
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

			AssertErrorMessageWhenStandupNotRegistered := func() func() {
				return func() {
					delete(store.StoreMap, context.Channel.Id)

					result := whiteboard.ProcessCommand(command+" "+title, context)

					expectedResult := MessageCommandResult{Text: MISSING_STANDUP, Status: THUMBS_DOWN}

					Expect(result).To(Equal(expectedResult))
				}
			}

			Describe("faces", func() {
				BeforeEach(func() {
					command = "faces"
					commitString = "Create New Face"
					expectedEntryItemKind = "New face"
					title = "Nicholas Cage"
					author = context.User.Author
					entry := NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					entry.Id = "1"
					expectedEntryType = Face{Entry: entry}
					expectedResult = EntryCommandResult{
						Title:    "NEW FACE",
						Status:   THUMBS_UP,
						HelpText: "",
						Entry:    entry,
					}
				})

				It("contains a new face entry in the result", AssertContainsEntryInResult())

				It("stores the new face entry in the entry map", AssertEntryStoredInEntryMap())

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

				Context("when the standup is not registered", func() {
					It("returns an error message", AssertErrorMessageWhenStandupNotRegistered())
				})
			})

			Describe("helps", func() {
				BeforeEach(func() {
					command = "helps"
					commitString = "Create Item"
					expectedEntryItemKind = "Help"
					title = "Good wicker furniture shop recommendations?"
					author = context.User.Author
					entry := NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					entry.Id = "1"
					expectedEntryType = Help{Entry: entry}
					expectedResult = EntryCommandResult{
						Title:    "HELP",
						Status:   THUMBS_UP,
						HelpText: NEW_ENTRY_HELP_TEXT,
						Entry:    entry,
					}
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

				Context("when the standup is not registered", func() {
					It("returns an error message", AssertErrorMessageWhenStandupNotRegistered())
				})
			})

			Describe("interestings", func() {
				BeforeEach(func() {
					command = "interestings"
					commitString = "Create Item"
					expectedEntryItemKind = "Interesting"
					title = "Nicholas Cage did a remake of The Wicker Man!"
					author = context.User.Author
					entry := NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					entry.Id = "1"
					expectedEntryType = Interesting{Entry: entry}
					expectedResult = EntryCommandResult{
						Title:    "INTERESTING",
						Status:   THUMBS_UP,
						HelpText: NEW_ENTRY_HELP_TEXT,
						Entry:    entry,
					}
				})

				It("contains an interesting entry in the result", AssertContainsEntryInResult())

				It("stores the interesting entry in the entry map", AssertEntryStoredInEntryMap())

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

				Context("when the standup is not registered", func() {
					It("returns an error message", AssertErrorMessageWhenStandupNotRegistered())
				})
			})

			Describe("events", func() {
				BeforeEach(func() {
					command = "events"
					commitString = "Create Item"
					expectedEntryItemKind = "Event"
					title = "Movie Screening for The Wicker Man!"
					author = context.User.Author
					entry := NewEntry(clock, author, title, sydneyStandup, expectedEntryItemKind)
					entry.Id = "1"
					expectedEntryType = Event{Entry: entry}
					expectedResult = EntryCommandResult{
						Title:    "EVENT",
						Status:   THUMBS_UP,
						HelpText: NEW_ENTRY_HELP_TEXT,
						Entry:    entry,
					}
				})

				It("contains an event entry in the result", AssertContainsEntryInResult())

				It("stores the event entry in the entry map", AssertEntryStoredInEntryMap())

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

				Context("when the standup is not registered", func() {
					It("returns an error message", AssertErrorMessageWhenStandupNotRegistered())
				})
			})
		})

		Describe("updating entries", func() {
			var (
				command           string
				newValue          string
				expectedEntryType EntryType
				expectedResult    EntryCommandResult
			)

			AssertEntryUpdated := func() func() {
				return func() {
					result := whiteboard.ProcessCommand(command+" "+newValue, context)

					Expect(result).To(Equal(expectedResult))
				}
			}

			AssertEntrySaved := func() func() {
				return func() {
					whiteboard.ProcessCommand(command+" "+newValue, context)

					Expect(gateway.SaveEntryCalled).To(BeTrue())
					Expect(gateway.EntrySaved).To(Equal(expectedEntryType))
				}
			}

			AssertNoInputErrorMessage := func() func() {
				return func() {
					expectedResult := MessageCommandResult{Text: MISSING_INPUT, Status: THUMBS_DOWN}

					result := whiteboard.ProcessCommand(command+"     ", context)

					Expect(result).To(Equal(expectedResult))
				}
			}

			AssertNoEntryErrorMessage := func() func() {
				return func() {
					expectedResult := MessageCommandResult{Text: MISSING_ENTRY, Status: THUMBS_DOWN}

					delete(whiteboard.EntryMap, context.User.Username)

					result := whiteboard.ProcessCommand(command+" "+newValue, context)

					Expect(result).To(Equal(expectedResult))
				}
			}

			AssertSaveEntryFailureErrorMessage := func() func() {
				return func() {
					gateway.SetSaveEntryError()
					expectedResult := MessageCommandResult{Text: "Problem creating post."}

					result := whiteboard.ProcessCommand(command+" "+newValue, context)

					Expect(result).To(Equal(expectedResult))
				}
			}

			Describe("body", func() {
				BeforeEach(func() {
					command = "body"
					newValue = "Nicholas Cage did a remake of The Wicker Man!"

					originalEntryId := "abc123"
					originalEntryType := NewInteresting(clock, context.User.Author, "Wicker Man 2006", sydneyStandup)
					originalEntryType.GetEntry().Id = originalEntryId
					whiteboard.EntryMap[context.User.Username] = originalEntryType

					expectedEntryType = NewInteresting(clock, context.User.Author, "Wicker Man 2006", sydneyStandup)
					newEntry := expectedEntryType.GetEntry()
					newEntry.Id = originalEntryId
					newEntry.Body = newValue
					expectedResult = EntryCommandResult{
						Title:    "INTERESTING",
						Status:   THUMBS_UP,
						HelpText: "",
						Entry:    newEntry,
					}
				})

				It("updates the body on the entry", AssertEntryUpdated())

				It("updates whiteboard with the body", AssertEntrySaved())

				Context("when the body is empty", func() {
					It("returns an error message", AssertNoInputErrorMessage())
				})

				Context("no entry in store", func() {
					It("returns an error message", AssertNoEntryErrorMessage())
				})

				Context("when saving the entry fails", func() {
					It("returns an error message", AssertSaveEntryFailureErrorMessage())
				})

				Context("when entry type is a New Face", func() {
					It("returns an error", func() {
						originalEntryId := "abc123"
						originalEntryType := NewFace(clock, context.User.Author, "Nicholas Cage", sydneyStandup)
						originalEntryType.GetEntry().Id = originalEntryId
						whiteboard.EntryMap[context.User.Username] = originalEntryType

						expectedResult := MessageCommandResult{
							Text:   "Hey, new faces should not have a body!",
							Status: THUMBS_DOWN,
						}

						result := whiteboard.ProcessCommand("body And John Travolta!", context)

						Expect(result).To(Equal(expectedResult))
					})
				})
			})

			Describe("date", func() {
				BeforeEach(func() {
					command = "date"
					newValue = "3000-05-13"

					originalEntryId := "abc123"
					originalEntryType := NewInteresting(clock, context.User.Author, "Wicker Man 2006", sydneyStandup)
					originalEntryType.GetEntry().Id = originalEntryId
					whiteboard.EntryMap[context.User.Username] = originalEntryType

					expectedEntryType = NewInteresting(clock, context.User.Author, "Wicker Man 2006", sydneyStandup)
					newEntry := expectedEntryType.GetEntry()
					newEntry.Id = originalEntryId
					newEntry.Date = newValue
					expectedResult = EntryCommandResult{
						Title:    "INTERESTING",
						Status:   THUMBS_UP,
						HelpText: "",
						Entry:    newEntry,
					}
				})

				It("updates the date of an entry", AssertEntryUpdated())

				It("updates whiteboard with the new date", AssertEntrySaved())

				Context("no entry in store", func() {
					It("returns an error message", AssertNoEntryErrorMessage())
				})

				Context("when the date is empty", func() {
					It("returns an error message", AssertNoInputErrorMessage())
				})

				Context("when saving the entry fails", func() {
					It("returns an error message", AssertSaveEntryFailureErrorMessage())
				})

				Context("when the date format is wrong", func() {
					It("returns an error", func() {
						errorMsg := "Date not set, use YYYY-MM-DD as date format\n"
						expectedResult := MessageCommandResult{Text: errorMsg, Status: THUMBS_DOWN}

						result := whiteboard.ProcessCommand("date LOLWUT", context)

						Expect(result).To(Equal(expectedResult))
					})
				})
			})

			Describe("name", func() {
				BeforeEach(func() {
					command = "name"
					newValue = "Olivia Newton John"
					originalEntryId := "abc123"
					originalEntryType := NewFace(clock, context.User.Author, "Oliver Newton John", sydneyStandup)
					originalEntryType.GetEntry().Id = originalEntryId
					whiteboard.EntryMap[context.User.Username] = originalEntryType

					expectedEntryType = NewFace(clock, context.User.Author, newValue, sydneyStandup)
					newEntry := expectedEntryType.GetEntry()
					newEntry.Id = originalEntryId
					expectedResult = EntryCommandResult{
						Title:    "NEW FACE",
						Status:   THUMBS_UP,
						HelpText: "",
						Entry:    newEntry,
					}
				})

				It("updates the name on a new face", AssertEntryUpdated())

				It("updates whiteboard with the new name", AssertEntrySaved())

				Context("when the new name is the empty string", func() {
					It("returns an error message", AssertNoInputErrorMessage())
				})

				Context("no entry in store", func() {
					It("returns an error message", AssertNoEntryErrorMessage())
				})

				Context("when saving the entry fails", func() {
					It("returns an error message", AssertSaveEntryFailureErrorMessage())
				})
			})

			Describe("title", func() {
				BeforeEach(func() {
					command = "title"
					newValue = "Saturday Night Live"
					originalEntryId := "abc123"
					originalEntryType := NewEvent(clock, context.User.Author, "Saturday Night Fever", sydneyStandup)
					originalEntryType.GetEntry().Id = originalEntryId
					whiteboard.EntryMap[context.User.Username] = originalEntryType

					expectedEntryType = NewEvent(clock, context.User.Author, newValue, sydneyStandup)
					newEntry := expectedEntryType.GetEntry()
					newEntry.Id = originalEntryId
					expectedResult = EntryCommandResult{
						Title:    "EVENT",
						Status:   THUMBS_UP,
						HelpText: "",
						Entry:    newEntry,
					}
				})

				It("updates the title on an entry", AssertEntryUpdated())

				It("updates whiteboard with the new title", AssertEntrySaved())

				Context("when the new title is the empty string", func() {
					It("returns an error message", AssertNoInputErrorMessage())
				})

				Context("no entry in store", func() {
					It("returns an error message", AssertNoEntryErrorMessage())
				})
			})
		})

		Describe("present", func() {
			It("fetches standup items", func() {
				standupText := ">>>— — —\n \n \n \nNEW FACES\n\n\n \n \n \nHELPS\n\n\n \n \n \nINTERESTINGS\n\n*Interesting 1*\n[Alice]\n02 Jan 2015\n \n*Interesting 2*\n[Bob]\n12 Jan 2015\n \n \n \nEVENTS\n\n\n \n \n \n— — —\n:clap:"
				expectedResult := MessageCommandResult{Text: standupText}

				result := whiteboard.ProcessCommand("present", context)

				Expect(result).To(Equal(expectedResult))
			})

			Context("when a number of days is specified", func() {
				It("pass that number of days on to filter the items", func() {
					standupText := ">>>— — —\n \n \n \nNEW FACES\n\n\n \n \n \nHELPS\n\n\n \n \n \nINTERESTINGS\n\n*Interesting 1*\n[Alice]\n02 Jan 2015\n \n \n \nEVENTS\n\n\n \n \n \n— — —\n:clap:"
					expectedResult := MessageCommandResult{Text: standupText}

					result := whiteboard.ProcessCommand("present 5", context)

					Expect(result).To(Equal(expectedResult))
				})
			})

			Context("when a number of days does not parse", func() {
				It("pass that number of days on to filter the items", func() {
					standupText := ">>>— — —\n \n \n \nNEW FACES\n\n\n \n \n \nHELPS\n\n\n \n \n \nINTERESTINGS\n\n*Interesting 1*\n[Alice]\n02 Jan 2015\n \n*Interesting 2*\n[Bob]\n12 Jan 2015\n \n \n \nEVENTS\n\n\n \n \n \n— — —\n:clap:"
					expectedResult := MessageCommandResult{Text: standupText}

					result := whiteboard.ProcessCommand("present forever", context)

					Expect(result).To(Equal(expectedResult))
				})
			})

			Context("when the standup is not registered", func() {
				It("returns an error message", func() {
					delete(store.StoreMap, context.Channel.Id)

					result := whiteboard.ProcessCommand("present", context)

					expectedResult := MessageCommandResult{
						Text:   MISSING_STANDUP,
						Status: THUMBS_DOWN,
					}

					Expect(result).To(Equal(expectedResult))
				})
			})

			Context("when retrieving the standup items fails", func() {
				It("returns an error message", func() {
					gateway.SetGetStandupItemsError()

					result := whiteboard.ProcessCommand("present", context)

					expectedResult := MessageCommandResult{
						Text:   "Error retrieving standup items.",
						Status: THUMBS_DOWN,
					}

					Expect(result).To(Equal(expectedResult))
				})
			})
		})
	})

	It("maintains current entry for each user", func() {
		bob := SlackUser{Username: "bob", Author: "Bob Bobbins"}
		channel := SlackChannel{Id: "C456", Name: "sydney-standup"}
		bobContext := SlackContext{User: bob, Channel: channel}

		andrewEntry := NewEntry(clock, "Andrew Leung", "Andrew's Interesting", sydneyStandup, "Interesting")
		andrewEntry.Id = "1"
		andrewEntry.Body = "This is very interesting."
		expectedAndrewResult := EntryCommandResult{
			Title:  "INTERESTING",
			Status: THUMBS_UP,
			Entry:  andrewEntry,
		}

		bobEntry := NewEntry(clock, "Bob Bobbins", "Bob's Awesome Event", sydneyStandup, "Event")
		bobEntry.Id = "1"
		expectedBobResult := EntryCommandResult{
			Title:  "EVENT",
			Status: THUMBS_UP,
			Entry:  bobEntry,
		}

		whiteboard.ProcessCommand("i Andrew's Interesting", context)
		whiteboard.ProcessCommand("e Bob's Event", bobContext)
		bobResult := whiteboard.ProcessCommand("t Bob's Awesome Event", bobContext)
		andrewResult := whiteboard.ProcessCommand("b This is very interesting.", context)

		Expect(andrewResult).To(Equal(expectedAndrewResult))
		Expect(bobResult).To(Equal(expectedBobResult))
	})
})
