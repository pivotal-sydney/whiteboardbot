package app

import (
	"github.com/nlopes/slack"
	"strings"
	"fmt"
)

type SlackClient interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}

func ParseMessageEvent(client SlackClient, ev *slack.MessageEvent) (username string, message string) {
	if strings.HasPrefix(ev.Text, "wb ") {
		user, err := client.GetUserInfo(ev.User)
		if err != nil {
			fmt.Printf("%v, %v", ev.User, err)
			return
		}
		message = strings.Join([]string{user.Name, ev.Text[3:]}, " ")
		client.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
		fmt.Printf("Posting message: %v", message)
		username = user.Name
	}
	return
}