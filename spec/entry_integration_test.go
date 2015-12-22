package spec_test

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"github.com/xtreme-andleung/whiteboardbot/model"
)

var _ = Describe("Entry Integration", func() {
	var (
		slackClient spec.MockSlackClient
		clock spec.MockClock
		restClient spec.MockRestClient
		whiteboard WhiteboardApp
		registrationEvent, usageEvent, newInterestingEvent, newEventEvent, newHelpEvent,
		newFaceEventTitleEvent, newInterestingWithTitleEvent, newHelpEventTitleEvent, newEventEventWithTitleEvent,
		setTitleEvent, setDateEvent, setBodyEvent MessageEvent
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		whiteboard = WhiteboardApp{SlackClient: &slackClient, Clock: clock, RestClient: &restClient, Store: &spec.MockStore{}, EntryMap: make(map[string]model.EntryType)}

		usageEvent = CreateMessageEvent("wb ?")
		newInterestingEvent = CreateMessageEvent("wb Intere")
		newEventEvent = CreateMessageEvent("wb Ev")
		newHelpEvent = CreateMessageEvent("wb hEl")
		newInterestingWithTitleEvent = CreateMessageEvent("wb\nint \n \n   something interesting")
		newEventEventWithTitleEvent = CreateMessageEvent("wb e\t\t\t\t\n          some event")
		newHelpEventTitleEvent = CreateMessageEvent("wb h some help")
		newFaceEventTitleEvent = CreateMessageEvent("wb f some face")
		setTitleEvent = CreateMessageEvent("Wb tI something interesting")
		setDateEvent = CreateMessageEvent("Wb dA 2015-12-01")
		setBodyEvent = CreateMessageEvent("wB Bod more info")

		registrationEvent = CreateMessageEvent("wb r 1")
		whiteboard.ParseMessageEvent(&registrationEvent)
	})

	Describe("with interesting keyword", func() {
		It("should begin creating a new interesting entry", func() {
			whiteboard.ParseMessageEvent(&newInterestingEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Interesting{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(BeEmpty())
			Expect(slackClient.Status).To(BeEmpty())
		})
	})

	Describe("with interesting keyword and title", func() {
		It("should create a new interesting entry with title", func() {
			whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Interesting{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(Equal("something interesting"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP))
		})
	})

	Describe("with help keyword and title", func() {
		It("should create a new help entry with title", func() {
			whiteboard.ParseMessageEvent(&newHelpEventTitleEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Help{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(Equal("some help"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP))
		})
	})

	Describe("with event keyword and title", func() {
		It("should create a new event entry with title", func() {
			whiteboard.ParseMessageEvent(&newEventEventWithTitleEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Event{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(Equal("some event"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP))
		})
	})

	Describe("with face keyword and title", func() {
		It("should create a new face entry with title", func() {
			whiteboard.ParseMessageEvent(&newFaceEventTitleEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Face{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(Equal("some face"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP))
		})
	})

	Describe("with event keyword", func() {
		It("should begin creating a new event entry", func() {
			whiteboard.ParseMessageEvent(&newEventEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Event{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(BeEmpty())
			Expect(slackClient.Status).To(BeEmpty())
		})
	})

	Describe("with help keyword", func() {
		It("should begin creating a new help entry", func() {
			whiteboard.ParseMessageEvent(&newHelpEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Help{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(BeEmpty())
			Expect(slackClient.Status).To(BeEmpty())
		})
	})

	Context("setting a title detail", func() {
		Describe("with an interesting entry started", func() {
			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newInterestingEvent)
			})
			Describe("with correct keyword", func() {
				It("should set the title of the entry", func() {
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(slackClient.EntryType.GetEntry().Title).To(Equal("something interesting"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP))
				})
				It("should post interesting entry to whiteboard since all mandatory fields are set", func() {
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					Expect(restClient.Request.Commit).To(Equal("Create Item"))
					Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
					Expect(restClient.Request.Item.StandupId).To(Equal(1))
				})
				It("should update existing interesting entry in the whiteboard ", func() {
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
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
					Expect(slackClient.Status).To(Equal(THUMBS_UP))
				})
				It("should not update existing interesting entry in the whiteboard when incorrect keyword", func() {
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					setTitleEvent.Text = "wb invalid"
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
				})
			})
			Describe("with non-keyword", func() {
				It("should respond with default", func() {
					setTitleEvent.Text = "wb titleSomethingWrong"
					whiteboard.ParseMessageEvent(&setTitleEvent)
					Expect(slackClient.Message).To(Equal("aleung no you titleSomethingWrong"))
				})
			})
			Describe("with question mark", func() {
				It("should respond with usage screen", func() {
					whiteboard.ParseMessageEvent(&usageEvent)
					Expect(slackClient.Message).Should(HavePrefix("*Usage*:\n    `wb [command] [text...]`"))
				})
			})
		})
		Describe("with no entry started", func() {
			It("should give a hint on how to start entry", func() {
				whiteboard.ParseMessageEvent(&setTitleEvent)
				Expect(slackClient.Message).To(Equal("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event]` first!"))
			})
		})
	})
	Context("setting a date detail", func() {
		Describe("with an interesting entry started", func() {

			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newInterestingEvent)
			})

			Describe("with correct keyword", func() {
				It("should set the date of the entry and respond with interesting string", func() {
					whiteboard.ParseMessageEvent(&setDateEvent)
					Expect(slackClient.EntryType.GetEntry().Date).To(Equal("2015-12-01"))
					Expect(slackClient.Status).To(BeEmpty())
				})
				It("should not set invalid date and respond with help message", func() {
					setDateEvent.Text = "wb date 12/01/2015"
					whiteboard.ParseMessageEvent(&setDateEvent)
					Expect(slackClient.EntryType.GetEntry().Date).To(Equal("2015-01-02"))
					Expect(slackClient.Status).To(Equal("\nDate not set, use YYYY-MM-DD as date format"))
				})
			})
		})

		Describe("with no entry started", func() {
			It("should give a hint on how to start entry", func() {
				whiteboard.ParseMessageEvent(&setDateEvent)
				Expect(slackClient.Message).To(Equal("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event]` first!"))
			})
		})
	})
	Context("setting a body detail", func() {
		Describe("with an interesting entry started", func() {
			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newInterestingEvent)
			})

			Describe("with correct keyword", func() {
				It("should set the body of the entry and respond with interesting string", func() {
					whiteboard.ParseMessageEvent(&setBodyEvent)
					Expect(slackClient.EntryType.GetEntry().Body).To(Equal("more info"))
					Expect(slackClient.Status).To(BeEmpty())
				})
			})
		})
		Describe("with no entry started", func() {
			It("should give a hint on how to start entry", func() {
				whiteboard.ParseMessageEvent(&setBodyEvent)
				Expect(slackClient.Message).To(Equal("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event]` first!"))
			})
		})
	})

	Context("with multiple users", func() {
		var (
			newEventAndrew, newEventDariusz,
			setNameAndrew, setNameDariusz MessageEvent
		)
		BeforeEach(func() {
			newEventAndrew = CreateMessageEventWithUser("wb f", "aleung")
			newEventDariusz = CreateMessageEventWithUser("wb i", "dlorenc")
			setNameAndrew = CreateMessageEventWithUser("wb n Andrew Leung", "aleung")
			setNameDariusz = CreateMessageEventWithUser("wb t Dariusz Lorenc", "dlorenc")
		})
		Describe("sending commands", func() {
			It("should create entries uniquely to each user", func() {
				whiteboard.ParseMessageEvent(&newEventAndrew)
				whiteboard.ParseMessageEvent(&newEventDariusz)
				whiteboard.ParseMessageEvent(&setNameAndrew)
				Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Face{}))
				whiteboard.ParseMessageEvent(&setNameDariusz)
				Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Interesting{}))
			})
		})
	})

	Context("posting to another standup ID", func() {
		BeforeEach(func() {
			registrationEvent.Text = "wb r 123"
			whiteboard.ParseMessageEvent(&registrationEvent)
		})

		Describe("when channel registered with another standup ID", func() {
			It("should post entry with correct standup ID", func() {
				whiteboard.ParseMessageEvent(&newInterestingWithTitleEvent)
				Expect(restClient.Request.Item.StandupId).To(Equal(123))
			})
		})
	})

})

func CreateMessageEvent(text string) (event MessageEvent) {
	return CreateMessageEventWithUser(text, "aleung")
}

func CreateMessageEventWithUser(text string, user string) (event MessageEvent) {
	event = MessageEvent{Msg: Msg{Text: text, User: user, Channel: "whiteboard-sydney"}}
	return
}
