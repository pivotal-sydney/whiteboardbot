package app

import (
	"fmt"
)

type QuietWhiteboard interface {
	HandleInput(string) Response
}

type QuietWhiteboardApp struct {
	RestClient RestClient
	Store      Store
	CommandMap map[string]func(input string) Response
}

type Response struct {
	Text string `json:"text"`
}

func NewQuietWhiteboard(restClient RestClient, store Store) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{RestClient: restClient}
	whiteboard.Store = store
	whiteboard.CommandMap = make(map[string]func(input string) Response)
	whiteboard.init()
	return
}

func (whiteboard QuietWhiteboardApp) init() {
	whiteboard.registerCommand("?", whiteboard.handleUsageCommand)
	whiteboard.registerCommand("register", whiteboard.handleRegistrationCommand)
}

func (whiteboard QuietWhiteboardApp) HandleInput(input string) Response {
	command, input := readNextCommand(input)
	return whiteboard.handleCommand(command, input)
}

func (whiteboard QuietWhiteboardApp) registerCommand(
	command string,
	callback func(input string) Response) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string) Response {
	return Response{USAGE}
}

func (whiteboard QuietWhiteboardApp) handleRegistrationCommand(standupId string) Response {
	standup, ok := whiteboard.RestClient.GetStandup(standupId)
	if !ok {
		return Response{"Standup not found!"}
	}

	whiteboard.Store.SetStandup(standupId, standup)

	text := fmt.Sprintf("Standup %v has been registered! You can now start creating Whiteboard entries!", standup.Title)

	return Response{text}
}

func (whiteboard QuietWhiteboardApp) handleCommand(command, input string) Response {
	for key, _ := range whiteboard.CommandMap {
		if matches(command, key) {
			callback := whiteboard.CommandMap[key]
			return callback(input)
		}
	}

	return Response{"Ooops"}
}
