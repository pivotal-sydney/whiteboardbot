package app

import (
	"github.com/nlopes/slack"
	"strings"
	"fmt"
	"github.com/xtreme-andleung/whiteboardbot/entry"
)

type SlackClient interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}

func ParseMessageEvent(client SlackClient, clock entry.Clock, ev *slack.MessageEvent) (username string, message string) {
	if strings.HasPrefix(ev.Text, "wb ") {
		user, err := client.GetUserInfo(ev.User)
		if err != nil {
			fmt.Printf("%v, %v", ev.User, err)
			return
		}
		username = user.Name
		message = ev.Text[3:]
		if strings.HasPrefix(message, "faces") {
			message = entry.NewFace(clock).String()
		} else {
			message = strings.Join([]string{user.Name, "no you", message}, " ")
		}
		fmt.Printf("Posting message: %v", message)
		client.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
	}
	return
}