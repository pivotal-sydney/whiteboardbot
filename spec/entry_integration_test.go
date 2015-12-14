package spec_test

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Entry Integration", func() {
	var (
		slackClient spec.MockSlackClient
		clock       spec.MockClock
		restClient  spec.MockRestClient

		usageEvent, newInterestingEvent, newEventEvent, newHelpEvent,
		newFaceEventTitleEvent, newInterestingWithTitleEvent, newHelpEventTitleEvent, newEventEventWithTitleEvent,
		setTitleEvent, setDateEvent, setBodyEvent MessageEvent
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		usageEvent = createMessageEvent("wb ?")
		newInterestingEvent = createMessageEvent("wb Intere")
		newEventEvent = createMessageEvent("wb Ev")
		newHelpEvent = createMessageEvent("wb hEl")
		newInterestingWithTitleEvent = createMessageEvent("wb int something interesting")
		newEventEventWithTitleEvent = createMessageEvent("wb e some event")
		newHelpEventTitleEvent = createMessageEvent("wb h some help")
		newFaceEventTitleEvent = createMessageEvent("wb f some face")
		setTitleEvent = createMessageEvent("Wb tI something interesting")
		setDateEvent = createMessageEvent("Wb dA 2015-12-01")
		setBodyEvent = createMessageEvent("wB Bod more info")
	})

	Describe("with interesting keyword", func() {
		It("should begin creating a new interesting entry and respond with interesting string", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingEvent)
			Expect(slackClient.Message).To(Equal("interestings\n  *title: \n  body: \n  date: 2015-01-02"))
		})
	})

	Describe("with interesting keyword and title", func() {
		It("should create a new interesting entry with title and respond with string", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingWithTitleEvent)
			Expect(slackClient.Message).To(Equal("interestings\n  *title: something interesting\n  body: \n  date: 2015-01-02\nitem created"))
		})
	})

	Describe("with help keyword and title", func() {
		It("should create a new help entry with title and respond with string", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &newHelpEventTitleEvent)
			Expect(slackClient.Message).To(Equal("helps\n  *title: some help\n  body: \n  date: 2015-01-02\nitem created"))
		})
	})

	Describe("with event keyword and title", func() {
		It("should create a new event entry with title and respond with string", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &newEventEventWithTitleEvent)
			Expect(slackClient.Message).To(Equal("events\n  *title: some event\n  body: \n  date: 2015-01-02\nitem created"))
		})
	})

	Describe("with face keyword and title", func() {
		It("should create a new face entry with title and respond with string", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &newFaceEventTitleEvent)
			Expect(slackClient.Message).To(Equal("faces\n  *name: some face\n  date: 2015-01-02\nitem created"))
		})
	})

	Describe("with event keyword", func() {
		It("should begin creating a new event entry and respond with event string", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &newEventEvent)
			Expect(slackClient.Message).To(Equal("events\n  *title: \n  body: \n  date: 2015-01-02"))
		})
	})

	Describe("with help keyword", func() {
		It("should begin creating a new help entry and respond with help string", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &newHelpEvent)
			Expect(slackClient.Message).To(Equal("helps\n  *title: \n  body: \n  date: 2015-01-02"))
		})
	})

	Context("setting a title detail", func() {
		Describe("with an interesting entry started", func() {
			BeforeEach(func() {
				ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingEvent)
			})
			Describe("with correct keyword", func() {
				It("should set the title of the entry and respond with interesting string", func() {
					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
					Expect(slackClient.Message).Should(HavePrefix("interestings\n  *title: something interesting\n  body: \n  date: 2015-01-02"))
				})
				It("should post interesting entry to whiteboard since all mandatory fields are set", func() {
					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					Expect(restClient.Request.Commit).To(Equal("Create Item"))
					Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
					Expect(slackClient.Message).Should(HaveSuffix("item created"))
				})
				It("should update existing interesting entry in the whiteboard ", func() {
					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					setTitleEvent.Text = "wb title updated title"
					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
					Expect(restClient.Request.Method).To(Equal("patch"))
					Expect(restClient.Request.Commit).To(Equal("Update Item"))
					Expect(restClient.Request.Item.Title).To(Equal("updated title"))
					Expect(restClient.Request.Item.Description).To(BeEmpty())
					Expect(restClient.Request.Item.Author).To(Equal("Andrew Leung"))
					Expect(restClient.Request.Id).To(Equal("1"))
					Expect(slackClient.Message).Should(HaveSuffix("item updated"))
				})
				It("should not update existing interesting entry in the whiteboard when incorrect keyword", func() {
					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					setTitleEvent.Text = "wb invalid"
					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					Expect(slackClient.Message).ShouldNot(HaveSuffix("item updated"))
				})
			})
			Describe("with non-keyword", func() {
				It("should respond with default", func() {
					setTitleEvent.Text = "wb titleSomethingWrong"
					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
					Expect(slackClient.Message).To(Equal("aleung no you titleSomethingWrong"))
				})
			})
			Describe("with question mark", func() {
				It("should respond with usage screen", func() {
					ParseMessageEvent(&slackClient, &restClient, clock, &usageEvent)
					Expect(slackClient.Message).Should(HavePrefix("*Usage*:\n    `wb [command] [text...]`"))
				})
			})
		})
	})
	Context("setting a date detail", func() {
		Describe("with an interesting entry started", func() {

			BeforeEach(func() {
				ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingEvent)
			})

			Describe("with correct keyword", func() {
				It("should set the date of the entry and respond with interesting string", func() {
					ParseMessageEvent(&slackClient, &restClient, clock, &setDateEvent)
					Expect(slackClient.Message).To(Equal("interestings\n  *title: \n  body: \n  date: 2015-12-01"))
				})
				It("should not set invalid date and respond with help message", func() {
					setDateEvent.Text = "wb date 12/01/2015"
					ParseMessageEvent(&slackClient, &restClient, clock, &setDateEvent)
					Expect(slackClient.Message).To(Equal("interestings\n  *title: \n  body: \n  date: 2015-01-02\nDate not set, use YYYY-MM-DD as date format"))
				})
			})
			Describe("with incorrect keyword", func() {
				It("should respond with default", func() {
					setDateEvent.Text = "wb date2015-12-01"
					ParseMessageEvent(&slackClient, &restClient, clock, &setDateEvent)
					Expect(slackClient.Message).To(Equal("aleung no you date2015-12-01"))
				})
			})
		})
	})
	Context("setting a body detail", func() {
		Describe("with an interesting entry started", func() {
			BeforeEach(func() {
				ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingEvent)
			})

			Describe("with correct keyword", func() {
				It("should set the body of the entry and respond with interesting string", func() {
					ParseMessageEvent(&slackClient, &restClient, clock, &setBodyEvent)
					Expect(slackClient.Message).To(Equal("interestings\n  *title: \n  body: more info\n  date: 2015-01-02"))
				})
			})
		})
	})

	Context("with multiple users", func() {
		var (
			newEventAndrew, newEventDariusz,
			setNameAndrew, setNameDariusz MessageEvent
		)
		BeforeEach(func() {
			newEventAndrew = createMessageEventWithUser("wb f", "aleung")
			newEventDariusz = createMessageEventWithUser("wb f", "dlorenc")
			setNameAndrew = createMessageEventWithUser("wb n Andrew Leung", "aleung")
			setNameDariusz = createMessageEventWithUser("wb n Dariusz Lorenc", "dlorenc")
		})
		Describe("sending commands", func() {
			It("should create entries uniquely to each user", func() {
				ParseMessageEvent(&slackClient, &restClient, clock, &newEventAndrew)
				ParseMessageEvent(&slackClient, &restClient, clock, &newEventDariusz)
				ParseMessageEvent(&slackClient, &restClient, clock, &setNameAndrew)
				Expect(slackClient.Message).To(Equal("faces\n  *name: Andrew Leung\n  date: 2015-01-02\nitem created"))
				ParseMessageEvent(&slackClient, &restClient, clock, &setNameDariusz)
				Expect(slackClient.Message).To(Equal("faces\n  *name: Dariusz Lorenc\n  date: 2015-01-02\nitem created"))
			})
		})
	})

})

func createMessageEvent(text string) (event MessageEvent) {
	return createMessageEventWithUser(text, "aleung")
}

func createMessageEventWithUser(text string, user string) (event MessageEvent) {
	event = MessageEvent{Msg: Msg{Text: text, User: user}}
	return
}
