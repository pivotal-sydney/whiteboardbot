package app

import "github.com/nlopes/slack"

type SlackWrapper interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}
