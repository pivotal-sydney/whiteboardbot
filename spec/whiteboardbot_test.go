package spec_test

import (
	"github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Whiteboardbot", func() {
	var (
		slackClient       spec.MockSlackClient
		clock             spec.MockClock
		restClient        spec.MockRestClient
		whiteboard 		  WhiteboardApp

		helloWorldEvent, randomEvent, registrationEvent slack.MessageEvent
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		whiteboard = WhiteboardApp{SlackClient: &slackClient, Clock: clock, RestClient: &restClient, Store: &spec.MockStore{}}

		helloWorldEvent = slack.MessageEvent{}
		helloWorldEvent.Text = "wb hello world"
		helloWorldEvent.Channel = "whiteboard-sydney"

		randomEvent = slack.MessageEvent{}
		randomEvent.Text = "wbsome other text"
		randomEvent.Channel = "whiteboard-sydney"

		registrationEvent = slack.MessageEvent{}
		registrationEvent.Text = "wb r 1"
		registrationEvent.Channel = "whiteboard-sydney"
	})

	Context("when receiving a MessageEvent", func() {
		Describe("with text containing keywords", func() {
			It("should post a message with text", func() {
				whiteboard.ParseMessageEvent(&registrationEvent)
				whiteboard.ParseMessageEvent(&helloWorldEvent)
				Expect(slackClient.PostMessageCalled).To(Equal(true))
				Expect(slackClient.Message).To(Equal("aleung no you hello world"))
			})
		})
		Describe("with text not containing keywords", func() {
			It("should ignore the event", func() {
				whiteboard.ParseMessageEvent(&randomEvent)
				Expect(slackClient.PostMessageCalled).To(Equal(false))
				Expect(slackClient.Message).To(BeEmpty())
			})
		})
	})
})
