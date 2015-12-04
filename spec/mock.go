package spec

import (
	"github.com/nlopes/slack"
)

type MockSlackClient struct {
	postMessageCalled bool
}

func (client *MockSlackClient) GetPostMessageCalled() (bool) {
	return client.postMessageCalled
}

func (client *MockSlackClient) PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error) {
	client.postMessageCalled = true
	return "channel", "timestamp", nil
}

func (client *MockSlackClient) GetUserInfo(user string) (*slack.User, error) {
	User := slack.User{}
	User.Name = "aleung"
	return &User, nil
}
