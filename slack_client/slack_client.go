package slack_client

import "github.com/nlopes/slack"

type SlackClient interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}
