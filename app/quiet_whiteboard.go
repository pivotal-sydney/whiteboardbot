package app

import (
	"fmt"
)

type QuietWhiteboard interface {
	ProcessCommand(string, SlackContext) CommandResult
}

type QuietWhiteboardApp struct {
	RestClient RestClient
	Store      Store
	CommandMap map[string]func(input string) CommandResult
}

type CommandResult struct {
	Text string
}

func NewQuietWhiteboard(restClient RestClient, store Store) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{RestClient: restClient}
	whiteboard.Store = store
	whiteboard.CommandMap = make(map[string]func(input string) CommandResult)
	whiteboard.init()
	return
}

func (whiteboard QuietWhiteboardApp) init() {
	whiteboard.registerCommand("?", whiteboard.handleUsageCommand)
	whiteboard.registerCommand("register", whiteboard.handleRegistrationCommand)
}

func (whiteboard QuietWhiteboardApp) ProcessCommand(input string, context SlackContext) CommandResult {
	command, input := readNextCommand(input)
	return whiteboard.handleCommand(command, input)
}

func (whiteboard QuietWhiteboardApp) registerCommand(
	command string,
	callback func(input string) CommandResult) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string) CommandResult {
	return CommandResult{USAGE}
}

func (whiteboard QuietWhiteboardApp) handleRegistrationCommand(standupId string) CommandResult {
	standup, ok := whiteboard.RestClient.GetStandup(standupId)
	if !ok {
		return CommandResult{"Standup not found!"}
	}

	whiteboard.Store.SetStandup(standupId, standup)

	text := fmt.Sprintf("Standup %v has been registered! You can now start creating Whiteboard entries!", standup.Title)

	return CommandResult{text}
}

func (whiteboard QuietWhiteboardApp) handleCommand(command, input string) CommandResult {
	for key := range whiteboard.CommandMap {
		if matches(command, key) {
			callback := whiteboard.CommandMap[key]
			return callback(input)
		}
	}

	return CommandResult{"Ooops"}
}
