package spec

import (
	. "github.com/nlopes/slack"
	"github.com/pivotal-sydney/whiteboardbot/app"
	"github.com/pivotal-sydney/whiteboardbot/model"
	"strconv"
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
	restClient.SetStandup(model.Standup{Id: 1, TimeZone: "Australia/Sydney", Title: "Sydney"})
	store := MockStore{}
	whiteboard := app.NewWhiteboard(&slackClient, &restClient, clock, &store)
	return whiteboard
}

func createWhiteboardAndRegisterStandup(standupId int) app.WhiteboardApp {
	whiteboard := createWhiteboard()
	registerStandup(whiteboard, standupId)
	return whiteboard
}

func registerStandup(whiteboard app.WhiteboardApp, standupId int) {
	registrationEvent := createMessageEvent("wb r " + strconv.Itoa(standupId))
	whiteboard.ParseMessageEvent(&registrationEvent)
}
