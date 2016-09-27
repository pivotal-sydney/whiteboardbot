package spec

import (
	"github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/slack"
)

var _ = Describe("SlackBotServer", func() {
	var (
		mockWhiteBoard MockQuietWhiteboard
		server         SlackBotServer
		slackUser      SlackUser
		slackChannel   SlackChannel
	)

	BeforeEach(func() {
		slackUser = SlackUser{Username: "aleung", Author: "Andrew Leung", TimeZone: "Australia/Sydney"}
		slackChannel = SlackChannel{Id: "C456", Name: "sydney-standup"}

		mockSlackClient := MockSlackClient{}
		mockSlackClient.AddSlackUser("U123", slackUser)
		mockSlackClient.AddSlackChannel("C456", slackChannel)
		mockWhiteBoard = MockQuietWhiteboard{}
		server = SlackBotServer{Whiteboard: &mockWhiteBoard, SlackClient: &mockSlackClient}
	})

	Describe("ProcessMessage", func() {
		Context("when the message begins with wb", func() {
			It("invokes ProcessCommand on whiteboard", func() {
				expectedContext := SlackContext{User: slackUser, Channel: slackChannel}

				ev := slack.MessageEvent{
					Msg: slack.Msg{
						Channel: "C456",
						User:    "U123",
						Text:    "wb       make me a    sandwich",
					}}

				server.ProcessMessage(&ev)

				Expect(mockWhiteBoard.HandleInputCalled).To(BeTrue())
				Expect(mockWhiteBoard.HandleInputArgs.Text).To(Equal("make me a    sandwich"))
				Expect(mockWhiteBoard.HandleInputArgs.Context).To(Equal(expectedContext))
			})

			It("posts the command results to the slack channel", func() {

			})
		})

		Context("when the message does not begin with wb", func() {})
	})
})
