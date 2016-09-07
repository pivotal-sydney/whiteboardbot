package spec

import (
	. "github.com/benjamintanweihao/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

var _ = Describe("Entry Integration", func() {
	var (
		whiteboard WhiteboardApp
		slackClient *MockSlackClient
		restClient *MockRestClient

		newInterestingEvent, newEventEvent, newHelpEvent,
		newFaceEventTitleEvent, newInterestingWithTitleEvent, newHelpEventTitleEvent, newEventEventWithTitleEvent,
		setTitleEvent, setDateEvent, setBodyEvent MessageEvent
	)

	BeforeEach(func() {
		whiteboard = createWhiteboardAndRegisterStandup(1)
		slackClient = whiteboard.SlackClient.(*MockSlackClient)
		restClient = whiteboard.RestClient.(*MockRestClient)

		newInterestingEvent = createMessageEvent("wb Intere")
		newEventEvent = createMessageEvent("wb Ev")
		newHelpEvent = createMessageEvent("wb hEl")
		newInterestingWithTitleEvent = createMessageEvent("wb\nint \n \n   something interesting")
		newEventEventWithTitleEvent = createMessageEvent("wb e\t\t\t\t\n          some event")
		newHelpEventTitleEvent = createMessageEvent("wb h some help")
		newFaceEventTitleEvent = createMessageEvent("wb f some face")
		setTitleEvent = createMessageEvent("Wb tI something interesting")
		setDateEvent = createMessageEvent("Wb dA 2015-12-01")
		setBodyEvent = createMessageEvent("wB Bod more info")
	})

	Describe("with interesting keyword without title", func() {
		It("should respond missing title", func() {
			whiteboard.ParseMessageEvent(&newInterestingEvent)
			Expect(slackClient.Message).To(Equal("Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"))
			Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
		})
	})

	Describe("with event keyword without title", func() {
		It("should respond missing title", func() {
			whiteboard.ParseMessageEvent(&newEventEvent)
			Expect(slackClient.Message).To(Equal("Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"))
			Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
		})
	})

	Describe("with help keyword without title", func() {
		It("should respond missing title", func() {
			whiteboard.ParseMessageEvent(&newHelpEvent)
			Expect(slackClient.Message).To(Equal("Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"))
			Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
		})
	})

	Describe("with interesting keyword and title", func() {
		It("should create a new interesting entry with title", func() {
			whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			Expect(slackClient.Entry.ItemKind).To(Equal("Interesting"))
			Expect(slackClient.Entry.Title).To(Equal("something interesting"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP + "_Now go update the details. Need help?_ `wb ?`\n\nINTERESTING\n"))
			Expect(restClient.PostCalledCount).To(Equal(1))
			Expect(restClient.Request.Commit).To(Equal("Create Item"))
			Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
			Expect(restClient.Request.Item.StandupId).To(Equal(1))
		})
	})

	Describe("with help keyword and title", func() {
		It("should create a new help entry with title", func() {
			whiteboard.ParseMessageEvent(&newHelpEventTitleEvent)
			Expect(slackClient.Entry.ItemKind).To(Equal("Help"))
			Expect(slackClient.Entry.Title).To(Equal("some help"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP + "_Now go update the details. Need help?_ `wb ?`\n\nHELP\n"))
			Expect(restClient.PostCalledCount).To(Equal(1))
			Expect(restClient.Request.Commit).To(Equal("Create Item"))
			Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
			Expect(restClient.Request.Item.StandupId).To(Equal(1))
		})
	})

	Describe("with event keyword and title", func() {
		It("should create a new event entry with title", func() {
			whiteboard.ParseMessageEvent(&newEventEventWithTitleEvent)
			Expect(slackClient.Entry.ItemKind).To(Equal("Event"))
			Expect(slackClient.Entry.Title).To(Equal("some event"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP + "_Now go update the details. Need help?_ `wb ?`\n\nEVENT\n"))
			Expect(restClient.PostCalledCount).To(Equal(1))
			Expect(restClient.Request.Commit).To(Equal("Create Item"))
			Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
			Expect(restClient.Request.Item.StandupId).To(Equal(1))
		})
	})

	Describe("with face keyword and title", func() {
		It("should create a new face entry with title", func() {
			whiteboard.ParseMessageEvent(&newFaceEventTitleEvent)
			Expect(slackClient.Entry.ItemKind).To(Equal("New face"))
			Expect(slackClient.Entry.Title).To(Equal("some face"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP + "_Now go update the details. Need help?_ `wb ?`\n\nNEW FACE\n"))
			Expect(restClient.PostCalledCount).To(Equal(1))
			Expect(restClient.Request.Commit).To(Equal("Create New Face"))
			Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
			Expect(restClient.Request.Item.StandupId).To(Equal(1))
		})
	})

	Describe("with interesting keyword and title containing slack-escaped characters", func() {
		It("should create a new interesting entry with correct title", func() {
			newInterestingWithTitleEvent.Text = "wb i useful &amp; &lt;interesting&gt;"
			whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			Expect(slackClient.Entry.Title).To(Equal("useful &amp; &lt;interesting&gt;"))
			Expect(restClient.Request.Item.Title).To(Equal("useful & <interesting>"))
		})
	})

	Describe("with interesting keyword and title containing slack user IDs", func() {
		It("should create a new interesting entry with user names", func() {
			newInterestingWithTitleEvent.Text = "wb i <@UUserId> likes <@UUserId2>"
			whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			Expect(slackClient.Entry.Title).To(Equal("@user-name likes @user-name-two"))
			Expect(restClient.Request.Item.Title).To(Equal("@user-name likes @user-name-two"))
		})
	})

	Describe("with interesting keyword and title containing slack channel IDs", func() {
		It("should create a new interesting entry with channel names", func() {
			newInterestingWithTitleEvent.Text = "wb i <#CChannelId> has moved to <#CChannelId2>"
			whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			Expect(slackClient.Entry.Title).To(Equal("#channel-name has moved to #channel-name-two"))
			Expect(restClient.Request.Item.Title).To(Equal("#channel-name has moved to #channel-name-two"))
		})
	})

	Context("setting a title detail", func() {
		Describe("with an interesting entry started", func() {
			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			})

			Describe("with correct keyword", func() {
				It("should update existing interesting entry in the whiteboard", func() {
					setTitleEvent.Text = "wb title updated title"
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
					Expect(restClient.Request.Method).To(Equal("patch"))
					Expect(restClient.Request.Commit).To(Equal("Update Item"))
					Expect(restClient.Request.Item.Title).To(Equal("updated title"))
					Expect(restClient.Request.Item.Description).To(BeEmpty())
					Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
					Expect(restClient.Request.Item.StandupId).To(Equal(1))
					Expect(restClient.Request.Id).To(Equal("1"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP + "INTERESTING\n"))
				})

				It("should update interesting entry with unescaped characters in title", func() {
					setTitleEvent.Text = "wb t useful &amp; &lt;interesting&gt;"
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(slackClient.Entry.Title).To(Equal("useful &amp; &lt;interesting&gt;"))
					Expect(restClient.Request.Item.Title).To(Equal("useful & <interesting>"))
				})

				It("should not allow to change title to empty", func() {
					setTitleEvent.Text = "wb title "
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					Expect(slackClient.Message).To(Equal("Oi! The title/name can't be empty!"))
					Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
				})

				It("should not update existing interesting entry in the whiteboard when incorrect keyword", func() {
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
					setTitleEvent.Text = "wb invalid"
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
				})
			})

			Describe("with non-keyword", func() {
				It("should respond with default", func() {
					setTitleEvent.Text = "wb titleSomethingWrong"
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(slackClient.Message).To(Equal("aleung no you titleSomethingWrong"))
				})
			})
		})

		Describe("with no entry started", func() {
			It("should give a hint on how to start entry", func() {
				whiteboard.ParseMessageEvent(&setTitleEvent)
				Expect(slackClient.Message).To(Equal("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!"))
				Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
			})
		})
	})

	Context("setting a date detail", func() {
		Describe("with an interesting entry started", func() {
			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			})

			Describe("with correct keyword", func() {
				It("should set the date of the entry and respond with interesting string", func() {
					whiteboard.ParseMessageEvent(&setDateEvent)
					Expect(slackClient.Entry.Date).To(Equal("2015-12-01"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP + "INTERESTING\n"))
				})

				It("should not set invalid date and respond with help message", func() {
					setDateEvent.Text = "wb date 12/01/2015"
					whiteboard.ParseMessageEvent(&setDateEvent)
					Expect(slackClient.Entry.Date).To(Equal("2015-01-02"))
					Expect(slackClient.Status).To(Equal(THUMBS_DOWN + "Date not set, use YYYY-MM-DD as date format\n"))
				})
			})
		})

		Describe("with no entry started", func() {
			It("should give a hint on how to start entry", func() {
				whiteboard.ParseMessageEvent(&setDateEvent)
				Expect(slackClient.Message).To(Equal("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!"))
				Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
			})
		})
	})

	Context("setting a body detail", func() {
		Describe("with an interesting entry started", func() {
			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			})

			Describe("with correct keyword", func() {
				It("should set the body of the entry and respond with interesting string", func() {
					whiteboard.ParseMessageEvent(&setBodyEvent)
					Expect(slackClient.Entry.Body).To(Equal("more info"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP + "INTERESTING\n"))
				})

				It("should set the of the entry with unescaped title", func() {
					setBodyEvent.Text = "wb b useful &amp; &lt;interesting&gt;"
					whiteboard.ParseMessageEvent(&setBodyEvent)
					Expect(slackClient.Entry.Body).To(Equal("useful &amp; &lt;interesting&gt;"))
					Expect(restClient.Request.Item.Description).To(Equal("useful & <interesting>"))
				})
			})
		})
		Describe("with no entry started", func() {
			It("should give a hint on how to start entry", func() {
				whiteboard.ParseMessageEvent(&setBodyEvent)
				Expect(slackClient.Message).To(Equal("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!"))
				Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
			})
		})
	})

	Context("with multiple users", func() {
		var (
			newEventAndrew, newEventDariusz,
			setNameAndrew, setNameDariusz MessageEvent
		)

		BeforeEach(func() {
			newEventAndrew = createMessageEventWithUser("wb f face", "aleung")
			newEventDariusz = createMessageEventWithUser("wb i interesting", "dlorenc")
			setNameAndrew = createMessageEventWithUser("wb n Andrew Leung", "aleung")
			setNameDariusz = createMessageEventWithUser("wb t Dariusz Lorenc", "dlorenc")
		})

		Describe("sending commands", func() {
			It("should create entries uniquely to each user", func() {
				whiteboard.ParseMessageEvent(&newEventAndrew)
				whiteboard.ParseMessageEvent(&newEventDariusz)
				whiteboard.ParseMessageEvent(&setNameAndrew)
				Expect(slackClient.Entry.ItemKind).To(Equal("New face"))
				whiteboard.ParseMessageEvent(&setNameDariusz)
				Expect(slackClient.Entry.ItemKind).To(Equal("Interesting"))
			})
		})
	})

	Context("posting to another standup ID", func() {
		BeforeEach(func() {
			registerStandup(whiteboard, 123)
		})

		Describe("when channel registered with another standup ID", func() {
			It("should post entry with correct standup ID", func() {
				whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
				Expect(restClient.Request.Item.StandupId).To(Equal(123))
			})
		})
	})
})
