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
		setNameEvent slack.MessageEvent
		setDateEvent slack.MessageEvent
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

		setNameEvent = slack.MessageEvent{}
		setNameEvent.Text = "wb name Dariusz Lorenc"

		setDateEvent = slack.MessageEvent{}
		setDateEvent.Text = "wb date 2015-12-01"
	})

	Context("when receiving a MessageEvent", func() {
		Describe("with text containing keywords", func() {
			It("should post a message with username and text", func() {
				Username, Text := ParseMessageEvent(&client, clock, &helloWorldEvent)
				Expect(client.GetPostMessageCalled()).To(Equal(true))
				Expect(Username).To(Equal("aleung"))
				Expect(Text).To(Equal("aleung no you hello world"))

			})
		})
		Describe("with text not containing keywords", func() {
			It("should ignore the event", func() {
				Username, Text := ParseMessageEvent(&client, clock, &randomEvent)
				Expect(client.GetPostMessageCalled()).To(Equal(false))
				Expect(Username).To(BeEmpty())
				Expect(Text).To(BeEmpty())
			})
		})
		Describe("with faces keyword", func() {
			It("should begin creating a new face entry and respond with face string", func() {
				_, Text := ParseMessageEvent(&client, clock, &newFaceEvent)
				Expect(Text).To(Equal("faces\n  *name: \n  date: 2015-01-02"))
			})
		})

		Context("setting a name detail", func() {
			Describe("with a new face entry started", func() {
				BeforeEach(func() {
					ParseMessageEvent(&client, clock, &newFaceEvent)
				})
				Describe("with correct keyword", func() {
					It("should set the name of the entry and respond with face string", func() {
						_, Text := ParseMessageEvent(&client, clock, &setNameEvent)
						Expect(Text).To(Equal("faces\n  *name: Dariusz Lorenc\n  date: 2015-01-02"))
					})
				})
				Describe("with incorrect keyword", func() {
					It("should respond with default", func() {
						setNameEvent.Text = "wb nameSomethingWrong"
						_, Text := ParseMessageEvent(&client, clock, &setNameEvent)
						Expect(Text).To(Equal("aleung no you nameSomethingWrong"))
					})
				})
			})
		})
		Context("setting a date detail", func() {
			Describe("with a new face entry started", func() {
				BeforeEach(func() {
					ParseMessageEvent(&client, clock, &newFaceEvent)
				})
				Describe("with correct keyword", func() {
					It("should set the date of the entry and respond with face string", func() {
						_, Text := ParseMessageEvent(&client, clock, &setDateEvent)
						Expect(Text).To(Equal("faces\n  *name: \n  date: 2015-12-01"))
					})
					It("should not set invalid date and respond with help message", func() {
						setDateEvent.Text = "wb date 12/01/2015"
						_, Text := ParseMessageEvent(&client, clock, &setDateEvent)
						Expect(Text).To(Equal("faces\n  *name: \n  date: 2015-01-02\nDate not set, use YYYY-MM-DD as date format"))
					})
				})
				Describe("with incorrect keyword", func() {
					It("should respond with default", func() {
						setDateEvent.Text = "wb date2015-12-01"
						_, Text := ParseMessageEvent(&client, clock, &setDateEvent)
						Expect(Text).To(Equal("aleung no you date2015-12-01"))
					})
				})
			})
		})
	})
})
