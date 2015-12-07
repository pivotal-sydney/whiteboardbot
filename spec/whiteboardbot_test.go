package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Whiteboardbot", func() {
	var (
		helloWorldEvent slack.MessageEvent
		randomEvent slack.MessageEvent
		client spec.MockSlackClient
		clock spec.MockClock
		restClient spec.MockRestClient
	)

	BeforeEach(func() {
		client = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}

		helloWorldEvent = slack.MessageEvent{}
		helloWorldEvent.Text = "wb hello world"

		randomEvent = slack.MessageEvent{}
		randomEvent.Text = "wbsome other text"
	})

	Context("when receiving a MessageEvent", func() {
		Describe("with text containing keywords", func() {
			It("should post a message with username and text", func() {
				Username, Text := ParseMessageEvent(&client, &restClient, clock, &helloWorldEvent)
				Expect(client.PostMessageCalled).To(Equal(true))
				Expect(Username).To(Equal("aleung"))
				Expect(Text).To(Equal("aleung no you hello world"))
			})
		})
		Describe("with text not containing keywords", func() {
			It("should ignore the event", func() {
				Username, Text := ParseMessageEvent(&client, &restClient, clock, &randomEvent)
				Expect(client.PostMessageCalled).To(Equal(false))
				Expect(Username).To(BeEmpty())
				Expect(Text).To(BeEmpty())
			})
		})
	})
})
