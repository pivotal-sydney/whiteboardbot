package app

import (
	"fmt"
	. "github.com/pivotal-sydney/whiteboardbot/model"
)

type QuietWhiteboard interface {
	ProcessCommand(string, SlackContext) CommandResult
}

type CommandHandler func(input string, context SlackContext) CommandResult

type QuietWhiteboardApp struct {
	RestClient RestClient
	Store      Store
	CommandMap map[string]CommandHandler
	Clock      Clock
}

type CommandResult struct {
	Text  string
	Entry fmt.Stringer
}

func NewQuietWhiteboard(restClient RestClient, store Store, clock Clock) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{RestClient: restClient}
	whiteboard.Store = store
	whiteboard.Clock = clock
	whiteboard.CommandMap = make(map[string]CommandHandler)
	whiteboard.init()
	return
}

func (whiteboard QuietWhiteboardApp) init() {
	whiteboard.registerCommand("?", whiteboard.handleUsageCommand)
	whiteboard.registerCommand("register", whiteboard.handleRegistrationCommand)
	whiteboard.registerCommand("faces", whiteboard.handleFacesCommand)
}

func (whiteboard QuietWhiteboardApp) ProcessCommand(input string, context SlackContext) CommandResult {
	command, input := readNextCommand(input)
	return whiteboard.handleCommand(command, input, context)
}

func (whiteboard QuietWhiteboardApp) registerCommand(
	command string,
	callback CommandHandler) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string, _ SlackContext) CommandResult {
	return CommandResult{Text: USAGE}
}

func (whiteboard QuietWhiteboardApp) handleRegistrationCommand(standupId string, context SlackContext) CommandResult {
	standup, ok := whiteboard.RestClient.GetStandup(standupId)
	if !ok {
		return CommandResult{Text: "Standup not found!"}
	}

	whiteboard.Store.SetStandup(context.Channel.ChannelId, standup)

	text := fmt.Sprintf("Standup %v has been registered! You can now start creating Whiteboard entries!", standup.Title)

	return CommandResult{Text: text}
}

func (whiteboard QuietWhiteboardApp) handleFacesCommand(input string, context SlackContext) CommandResult {
	var entry fmt.Stringer

	if entry = resultIfEmptyTitle(input); entry == nil {
		standup, _ := whiteboard.Store.GetStandup(context.Channel.ChannelId)
		face := NewFace(whiteboard.Clock, context.User.Author, input, standup).(EntryType)
		entry = *face.GetEntry()
	}

	return CommandResult{Entry: entry}
}

func resultIfEmptyTitle(input string) fmt.Stringer {
	if len(input) == 0 {
		return InvalidEntry{Error: THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"}
	}
	return nil
}

func (whiteboard QuietWhiteboardApp) handleCommand(command, input string, context SlackContext) CommandResult {
	for key := range whiteboard.CommandMap {
		if matches(command, key) {
			callback := whiteboard.CommandMap[key]
			return callback(input, context)
		}
	}

	return CommandResult{Text: "Ooops"}
}
