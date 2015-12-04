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
		newFaceEvent slack.MessageEvent
		randomEvent slack.MessageEvent
		client spec.MockSlackClient
		clock spec.MockClock
	)

	BeforeEach(func() {
		client = spec.MockSlackClient{}
		helloWorldEvent = slack.MessageEvent{}
		helloWorldEvent.Text = "wb hello world"

		newFaceEvent = slack.MessageEvent{}
		newFaceEvent.Text = "wb faces"

		randomEvent = slack.MessageEvent{}
		randomEvent.Text = "wbsome other text"
		clock = spec.MockClock{}
	})

	Describe("when receiving a MessageEvent", func() {
		Context("with text containing keywords", func() {
			It("should post a message with username and text", func() {
				Username, Text := ParseMessageEvent(&client, clock, &helloWorldEvent)
				Expect(client.GetPostMessageCalled()).To(Equal(true))
				Expect(Username).To(Equal("aleung"))
				Expect(Text).To(Equal("aleung no you hello world"))

			})
		})
		Context("with text not containing keywords", func() {
			It("ignore the event", func() {
				Username, Text := ParseMessageEvent(&client, clock, &randomEvent)
				Expect(client.GetPostMessageCalled()).To(Equal(false))
				Expect(Username).To(BeEmpty())
				Expect(Text).To(BeEmpty())
			})
		})
		Context("with faces keyword", func() {
			It("should begin creating a new face entry and respond with face string", func() {
				_, Text := ParseMessageEvent(&client, clock, &newFaceEvent)
				Expect(Text).To(Equal("faces\n  *name: \n  date: 2015-01-02"))
			})
		})
	})
})
