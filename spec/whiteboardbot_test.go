package spec

import (
	"github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

var _ = Describe("Whiteboardbot", func() {
	var (
		whiteboard 		  WhiteboardApp
		slackClient       *MockSlackClient
		helloWorldEvent, randomEvent, registrationEvent slack.MessageEvent
	)

	BeforeEach(func() {
		whiteboard = createWhiteboard()
		slackClient = whiteboard.SlackClient.(*MockSlackClient)

		helloWorldEvent = createMessageEvent("wb hello world")
		randomEvent = createMessageEvent("wbsome other text")
		registrationEvent = createMessageEvent("wb r 1")
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