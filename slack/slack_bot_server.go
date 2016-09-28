package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

type SlackBotServer struct {
	SlackClient SlackClient
	Whiteboard  QuietWhiteboard
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
	server.SlackClient.PostMessage(result.String(), slackChannel.Id, THUMBS_UP)
}

func (server SlackBotServer) Run(rtm *slack.RTM) {
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				go server.ProcessMessage(ev)
			case *slack.InvalidAuthEvent:
				fmt.Println("Invalid credentials")
				break
			default:
			}
		}
	}
}
