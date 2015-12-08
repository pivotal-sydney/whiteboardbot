package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Interestings Integration", func() {
	var (
		newInterestingEvent slack.MessageEvent
		setTitleEvent slack.MessageEvent
		setDateEvent slack.MessageEvent
		slackClient spec.MockSlackClient
		clock spec.MockClock
		restClient spec.MockRestClient
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}

		newInterestingEvent = slack.MessageEvent{}
		newInterestingEvent.Text = "wb interestings"

		setTitleEvent = slack.MessageEvent{}
		setTitleEvent.Text = "wb title something interesting"

		setDateEvent = slack.MessageEvent{}
		setDateEvent.Text = "wb date 2015-12-01"
	})

	Describe("with interestings keyword", func() {
		It("should begin creating a new interesting entry and respond with interesting string", func() {
			_, Text := ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingEvent)
			Expect(Text).To(Equal("interestings\n  *title: \n  body: \n  date: 2015-01-02"))
		})
	})

//	Context("setting a title detail", func() {
//		Describe("with a new interesting entry started", func() {
//			BeforeEach(func() {
//				ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingEvent)
//			})
//			Describe("with correct keyword", func() {
//				It("should set the title of the entry and respond with interesting string", func() {
//					_, Text := ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
//					Expect(Text).Should(HavePrefix("interestings\n  *title: something interesting\n  body: \n  date: 2015-01-02"))
//				})
//				It("should post new interesting entry to whiteboard since all mandatory fields are set", func() {
//					_, message := ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
//					Expect(restClient.PostCalledCount).To(Equal(1))
//					Expect(restClient.Request.Commit).To(Equal("Create New Interesting"))
//					Expect(message).Should(HaveSuffix("new interesting created"))
//				})
//				It("should update existing interesting entry in the whiteboard ", func() {
//					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
//					Expect(restClient.PostCalledCount).To(Equal(1))
//					setTitleEvent.Text = "wb title updated title"
//					_, message := ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
//					Expect(restClient.PostCalledCount).To(Equal(2))
//					Expect(restClient.Request.Method).To(Equal("patch"))
//					Expect(restClient.Request.Commit).To(Equal("Update New Interesting"))
//					Expect(restClient.Request.Item.Title).To(Equal("updated title"))
//					Expect(restClient.Request.Id).To(Equal("1"))
//					Expect(message).Should(HaveSuffix("new interesting updated"))
//				})
//				It("should not update existing interesting entry in the whiteboard when incorrect keyword", func() {
//					ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
//					Expect(restClient.PostCalledCount).To(Equal(1))
//					setTitleEvent.Text = "wb invalid"
//					_, message := ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
//					Expect(restClient.PostCalledCount).To(Equal(1))
//					Expect(message).ShouldNot(HaveSuffix("new interesting updated"))
//				})
//			})
//			Describe("with incorrect keyword", func() {
//				It("should respond with default", func() {
//					setTitleEvent.Text = "wb titleSomethingWrong"
//					_, Text := ParseMessageEvent(&slackClient, &restClient, clock, &setTitleEvent)
//					Expect(Text).To(Equal("aleung no you titleSomethingWrong"))
//				})
//			})
//		})
//	})
//	Context("setting a date detail", func() {
//		Describe("with a new interesting entry started", func() {
//			BeforeEach(func() {
//				ParseMessageEvent(&slackClient, &restClient, clock, &newInterestingEvent)
//			})
//			Describe("with correct keyword", func() {
//				It("should set the date of the entry and respond with interesting string", func() {
//					_, Text := ParseMessageEvent(&slackClient, &restClient, clock, &setDateEvent)
//					Expect(Text).To(Equal("interestings\n  *title: \n  date: 2015-12-01"))
//				})
//				It("should not set invalid date and respond with help message", func() {
//					setDateEvent.Text = "wb date 12/01/2015"
//					_, Text := ParseMessageEvent(&slackClient, &restClient, clock, &setDateEvent)
//					Expect(Text).To(Equal("interestings\n  *title: \n  date: 2015-01-02\nDate not set, use YYYY-MM-DD as date format"))
//				})
//			})
//			Describe("with incorrect keyword", func() {
//				It("should respond with default", func() {
//					setDateEvent.Text = "wb date2015-12-01"
//					_, Text := ParseMessageEvent(&slackClient, &restClient, clock, &setDateEvent)
//					Expect(Text).To(Equal("aleung no you date2015-12-01"))
//				})
//			})
//		})
//	})
})


