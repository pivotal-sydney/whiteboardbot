package spec

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
)

var _ = Describe("Standup Registration", func() {
	var (
		whiteboard WhiteboardApp
		slackClient *MockSlackClient

		anythingEvent, registrationEvent MessageEvent
	)

	BeforeEach(func() {
		whiteboard = createWhiteboard()
		slackClient = whiteboard.SlackClient.(*MockSlackClient)

		anythingEvent = createMessageEvent("wb anything")
		registrationEvent = createMessageEvent("wb r 1")
	})

	Context("registering standup", func() {
		Describe("when standup has not been registered", func() {
			It("should ask for standup ID", func() {
				whiteboard.ParseMessageEvent(&anythingEvent)
				Expect(slackClient.Message).To(Equal("You haven't registered your standup yet. wb r <id> first!"))
				Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
			})
		})

		Describe("with an integer as standup id", func() {
			It("should respond registration successful", func() {
				whiteboard.ParseMessageEvent(&registrationEvent)
				Expect(slackClient.Message).To(Equal("Standup Sydney has been registered! You can now start creating Whiteboard entries!"))
				Expect(slackClient.Status).To(Equal(THUMBS_UP))
			})
		})
	})
})