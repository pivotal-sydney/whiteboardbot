package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/nlopes/slack"
	"whiteboardbot/interfaces"
)

var _ = Describe("Whiteboardbot", func() {
	var (
		helloWorldEvent slack.MessageEvent
		randomEvent slack.MessageEvent
		client interfaces.MockSlackClient
	)

	BeforeEach(func() {
		client = interfaces.MockSlackClient{}
		helloWorldEvent = slack.MessageEvent{}
		helloWorldEvent.Text = "wb hello world"

		randomEvent = slack.MessageEvent{}
		randomEvent.Text = "wbsome other text"
	})

	Describe("When receiving a MessageEvent", func() {
		Context("with text containing keywords", func() {
			It("should post a message with username and text", func() {
				Username, Text := interfaces.ParseMessageEvent(&client, &helloWorldEvent)
				Expect(client.GetPostMessageCalled()).To(Equal(true))
				Expect(Username).To(Equal("aleung"))
				Expect(Text).To(Equal("aleung hello world"))

			})
		})
		Context("with text not containing keywords", func() {
			It("ignore the event", func() {
				Username, Text := interfaces.ParseMessageEvent(&client, &randomEvent)
				Expect(client.GetPostMessageCalled()).To(Equal(false))
				Expect(Username).To(BeEmpty())
				Expect(Text).To(BeEmpty())
			})
		})
	})
})
