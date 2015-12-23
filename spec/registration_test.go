package spec_test

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Standup Registration", func() {
	var (
		slackClient spec.MockSlackClient
		clock spec.MockClock
		restClient spec.MockRestClient
		whiteboard WhiteboardApp

		event MessageEvent
		registrationEvent MessageEvent
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		whiteboard = WhiteboardApp{SlackClient: &slackClient, Clock: clock, RestClient: &restClient, Store: &spec.MockStore{}}

		event = CreateMessageEvent("wb anything")
		registrationEvent = CreateMessageEvent("wb r 1")
	})

	Context("registering standup", func() {
		Describe("when standup has not been registered", func() {
			It("should ask for standup ID", func() {
				whiteboard.ParseMessageEvent(&event)
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
