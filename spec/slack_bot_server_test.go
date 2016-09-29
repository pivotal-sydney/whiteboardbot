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
		mockWhiteBoard  MockWhiteboard
		mockSlackClient MockSlackClient
		server          SlackBotServer
		slackUser       SlackUser
		slackChannel    SlackChannel
	)

	BeforeEach(func() {
		slackUser = SlackUser{Username: "aleung", Author: "Andrew Leung", TimeZone: "Australia/Sydney"}
		slackChannel = SlackChannel{Id: "C456", Name: "sydney-standup"}

		mockSlackClient = MockSlackClient{}
		mockSlackClient.AddSlackUser("U123", slackUser)
		mockSlackClient.AddSlackChannel("C456", slackChannel)
		mockWhiteBoard = MockWhiteboard{}
		server = SlackBotServer{Whiteboard: &mockWhiteBoard, SlackClient: &mockSlackClient}
	})

	Describe("ProcessMessage", func() {
		Context("when the message begins with wb", func() {
			It("invokes ProcessCommand on whiteboard", func() {
				expectedContext := SlackContext{User: slackUser, Channel: slackChannel}
				messageEvent := makeMessageEvent("C456", "U123", "wb       make me a    sandwich")

				server.ProcessMessage(&messageEvent)

				Expect(mockWhiteBoard.HandleInputCalled).To(BeTrue())
				Expect(mockWhiteBoard.HandleInputArgs.Text).To(Equal("make me a    sandwich"))
				Expect(mockWhiteBoard.HandleInputArgs.Context).To(Equal(expectedContext))
			})

			It("posts the command results to the slack channel", func() {
				messageEvent := makeMessageEvent("C456", "U123", "wb       make me a    sandwich")
				server.ProcessMessage(&messageEvent)

				Expect(mockSlackClient.PostMessageCalled).To(BeTrue())
				Expect(mockSlackClient.Message).To(Equal("This is a mock message"))
				Expect(mockSlackClient.ChannelId).To(Equal("C456"))
			})

			It("replaces user IDs with usernames", func() {
				messageEvent := makeMessageEvent("C456", "U123", "wb i <@UUserId> likes <@UUserId2>")
				server.ProcessMessage(&messageEvent)
				Expect(mockWhiteBoard.HandleInputArgs.Text).To(Equal("i @user-name likes @user-name-two"))
			})

			It("replaces channel IDs with channel names", func() {
				messageEvent := makeMessageEvent("C456", "U123", "wb i <#CChannelId> has moved to <#CChannelId2>")
				server.ProcessMessage(&messageEvent)
				Expect(mockWhiteBoard.HandleInputArgs.Text).To(Equal("i #channel-name has moved to #channel-name-two"))
			})
		})

		Context("when the message does not begin with wb", func() {
			BeforeEach(func() {
				messageEvent := makeMessageEvent("C456", "U123", "       make me a    sandwich")
				server.ProcessMessage(&messageEvent)
			})

			It("does not invoke ProcessCommand on whiteboard", func() {
				Expect(mockWhiteBoard.HandleInputCalled).To(BeFalse())
			})

			It("does not post the entry", func() {
				Expect(mockSlackClient.PostMessageCalled).To(BeFalse())
			})

		})
	})
})

func makeMessageEvent(channel, user, text string) slack.MessageEvent {
	return slack.MessageEvent{
		Msg: slack.Msg{
			Channel: channel,
			User:    user,
			Text:    text,
		}}
}
