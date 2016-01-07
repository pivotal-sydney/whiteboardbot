package spec

import (
	. "github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/app"
)

func createMessageEvent(text string) MessageEvent {
	return createMessageEventWithUser(text, "aleung")
}

func createMessageEventWithUser(text string, user string) MessageEvent {
	return MessageEvent{Msg: Msg{Text: text, User: user, Channel: "whiteboard-sydney"}}
}

func createWhiteboard() app.WhiteboardApp {
	slackClient := MockSlackClient{}
	clock := MockClock{}
	restClient := MockRestClient{}
	store := MockStore{}
	return app.NewWhiteboard(&slackClient, &restClient, clock, &store)
}