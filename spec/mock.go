package spec

import (
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/model"
	"time"
)

type MockSlackClient struct {
	PostMessageCalled bool
	Message           string
}

func (client *MockSlackClient) PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error) {
	client.PostMessageCalled = true
	client.Message = text
	return "channel", "timestamp", nil
}

func (client *MockSlackClient) GetUserInfo(user string) (*slack.User, error) {
	slackUser := slack.User{}
	slackUser.Profile = slack.UserProfile{RealName: "Andrew Leung"}
	if user == "" {
		slackUser.Name = "aleung"
	} else {
		slackUser.Name = user
	}
	return &slackUser, nil
}

type MockClock struct{}

func (clock MockClock) Now() time.Time {
	return time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)
}

type MockRestClient struct {
	PostCalledCount int
	Request         model.WhiteboardRequest
}

func (client *MockRestClient) Post(request model.WhiteboardRequest) (itemId string, ok bool) {
	client.PostCalledCount++
	client.Request = request
	ok = true
	itemId = "1"
	return
}
