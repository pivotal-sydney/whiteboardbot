package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	"regexp"
	"strings"
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

	input = server.replaceIdsWithNames(input)

	result := server.Whiteboard.ProcessCommand(input, context)
	server.SlackClient.PostMessage(result.String(), slackChannel.Id)
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

func (server SlackBotServer) replaceIdsWithNames(input string) string {
	input = server.replaceUserIdsWithNames(input)
	input = server.replaceChannelIdsWithNames(input)
	return input
}

func (server SlackBotServer) replaceUserIdsWithNames(input string) string {
	re := regexp.MustCompile("<@([a-zA-Z0-9]+)>")
	userIds := re.FindAllStringSubmatch(input, -1)

	for _, id := range userIds {
		userId := id[1]
		slackUser := server.SlackClient.GetUserDetails(userId)
		userName := "@" + slackUser.Username
		input = strings.Replace(input, id[0], userName, -1)
	}

	return input
}

func (server SlackBotServer) replaceChannelIdsWithNames(input string) string {
	re := regexp.MustCompile("<#([a-zA-Z0-9]+)>")
	channelIds := re.FindAllStringSubmatch(input, -1)

	for _, id := range channelIds {
		channelId := id[1]
		slackChannel := server.SlackClient.GetChannelDetails(channelId)
		channelName := "#" + slackChannel.Name
		input = strings.Replace(input, id[0], channelName, -1)
	}

	return input
}
