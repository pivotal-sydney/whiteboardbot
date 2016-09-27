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
	command, input := ReadNextCommand(ev.Msg.Text)

	if command != "wb" {
		return
	}

	slackChannel := server.SlackClient.GetChannelDetails(ev.Msg.Channel)
	slackUser := server.SlackClient.GetUserDetails(ev.Msg.User)

	context := SlackContext{
		Channel: slackChannel,
		User:    slackUser,
	}

	result := server.Whiteboard.ProcessCommand(input, context)
	server.SlackClient.PostMessage(result.Entry.String(), slackChannel.Id, THUMBS_UP)

}
