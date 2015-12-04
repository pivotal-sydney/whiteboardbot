package interfaces

import (
	"github.com/nlopes/slack"
	"strings"
	"fmt"
)

type MockSlackClient struct {
	postMessageCalled bool
}

func (client MockSlackClient) GetPostMessageCalled() (bool) {
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

type SlackClient interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}

func ParseMessageEvent(client SlackClient, ev *slack.MessageEvent) (string, string) {
	message := ev.Text

	if strings.HasPrefix(message, "wb ") {
		message = message[3:len(message)]

		user, err := client.GetUserInfo(ev.User)

		if err != nil {
			fmt.Printf("%v, %v", ev.User, err)
			return "", ""
		}

		fmt.Printf("%v", user.Name)
		message = strings.Join([]string{user.Name, message}, " ")
		client.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
		return user.Name, message
	}
	return "", ""
}