package slack

import (
	"github.com/nlopes/slack"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

type SlackBotServer struct {
	Whiteboard  QuietWhiteboard
	SlackClient SlackClient
}

func (server SlackBotServer) ProcessMessage(ev *slack.MessageEvent) {
	_, input := ReadNextCommand(ev.Msg.Text)

	slackChannel := server.SlackClient.GetChannelDetails(ev.Msg.Channel)
	slackUser := server.SlackClient.GetUserDetails(ev.Msg.User)

	context := SlackContext{
		Channel: slackChannel,
		User:    slackUser,
	}

	server.Whiteboard.ProcessCommand(input, context)
}
