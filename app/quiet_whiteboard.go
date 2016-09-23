package app

import (
	"errors"
	"fmt"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"time"
)

type QuietWhiteboard interface {
	ProcessCommand(string, SlackContext) CommandResult
}

type CommandHandler func(input string, context SlackContext) CommandResult

type EntryFactory func(clock Clock, author, title string, standup Standup) EntryType

type QuietWhiteboardApp struct {
	Clock      Clock
	RestClient RestClient
	Repository StandupRepository
	Store      Store
	CommandMap map[string]CommandHandler
	EntryMap   map[string]EntryType
}

type CommandResult struct {
	Entry fmt.Stringer
}

func NewQuietWhiteboard(restClient RestClient, gateway StandupRepository, store Store, clock Clock) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{
		Clock:      clock,
		RestClient: restClient,
		Repository: gateway,
		Store:      store,
		CommandMap: make(map[string]CommandHandler),
		EntryMap:   make(map[string]EntryType),
	}
	whiteboard.init()
	return
}

func (whiteboard QuietWhiteboardApp) init() {
	whiteboard.registerCommand("?", whiteboard.handleUsageCommand)
	whiteboard.registerCommand("register", whiteboard.handleRegistrationCommand)
	whiteboard.registerCommand("faces", whiteboard.handleFacesCommand)
	whiteboard.registerCommand("helps", whiteboard.handleHelpsCommand)
	whiteboard.registerCommand("interestings", whiteboard.handleInterestingsCommand)
	whiteboard.registerCommand("events", whiteboard.handleEventsCommand)
	whiteboard.registerCommand("body", whiteboard.handleBodyCommand)
	whiteboard.registerCommand("date", whiteboard.handleDateCommand)
	whiteboard.registerCommand("name", whiteboard.handleUpdateCommand)
	whiteboard.registerCommand("title", whiteboard.handleUpdateCommand)
}

func (whiteboard QuietWhiteboardApp) ProcessCommand(input string, context SlackContext) CommandResult {
	command, input := readNextCommand(input)
	return whiteboard.handleCommand(command, input, context)
}

func (whiteboard QuietWhiteboardApp) handleCommand(command, input string, context SlackContext) CommandResult {
	for key := range whiteboard.CommandMap {
		if matches(command, key) {
			callback := whiteboard.CommandMap[key]
			return callback(input, context)
		}
	}

	return CommandResult{Entry: InvalidEntry{Error: "Ooops"}}
}

func (whiteboard QuietWhiteboardApp) registerCommand(command string, callback CommandHandler) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string, _ SlackContext) CommandResult {
	return CommandResult{Entry: TextEntry{Text: USAGE}}
}

func (whiteboard QuietWhiteboardApp) handleRegistrationCommand(standupId string, context SlackContext) (command CommandResult) {
	command = CommandResult{}

	standup, ok := whiteboard.RestClient.GetStandup(standupId)
	if !ok {
		command = CommandResult{Entry: InvalidEntry{Error: "Standup not found!"}}
		return
	}

	whiteboard.Store.SetStandup(context.Channel.ChannelId, standup)

	text := fmt.Sprintf("Standup %v has been registered! You can now start creating Whiteboard entries!", standup.Title)

	command.Entry = TextEntry{Text: text}
	return
}

func (whiteboard QuietWhiteboardApp) handleFacesCommand(input string, context SlackContext) CommandResult {
	return whiteboard.handleCreateCommand(input, context, NewFace)
}

func (whiteboard QuietWhiteboardApp) handleHelpsCommand(input string, context SlackContext) CommandResult {
	return whiteboard.handleCreateCommand(input, context, NewHelp)
}

func (whiteboard QuietWhiteboardApp) handleInterestingsCommand(input string, context SlackContext) CommandResult {
	return whiteboard.handleCreateCommand(input, context, NewInteresting)
}

func (whiteboard QuietWhiteboardApp) handleEventsCommand(input string, context SlackContext) CommandResult {
	return whiteboard.handleCreateCommand(input, context, NewEvent)
}

func (whiteboard QuietWhiteboardApp) handleBodyCommand(input string, context SlackContext) CommandResult {
	if err := handleEmptyInput(input); err != nil {
		return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
	}

	username := context.User.Username

	if entryType, ok := whiteboard.EntryMap[username]; ok {
		entry := entryType.GetEntry()

		if entry.ItemKind == "New face" {
			errorMsg := ":-1:\nHey, new faces should not have a body!"
			return CommandResult{InvalidEntry{Error: errorMsg}}
		}

		entryType.GetEntry().Body = input

		if _, err := whiteboard.Repository.SaveEntry(entryType); err != nil {
			return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
		}
		return CommandResult{Entry: entry}
	} else {
		return CommandResult{Entry: InvalidEntry{Error: MISSING_ENTRY}}
	}
}

func (whiteboard QuietWhiteboardApp) handleDateCommand(input string, context SlackContext) CommandResult {
	if err := handleEmptyInput(input); err != nil {
		return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
	}

	if parsedDate, err := time.Parse(DATE_FORMAT, input); err == nil {
		if entryType, ok := whiteboard.EntryMap[context.User.Username]; ok {
			entryType.GetEntry().Date = parsedDate.Format(DATE_FORMAT)

			if _, err := whiteboard.Repository.SaveEntry(entryType); err != nil {
				return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
			}

			return CommandResult{Entry: entryType.GetEntry()}
		} else {
			return CommandResult{Entry: InvalidEntry{Error: MISSING_ENTRY}}
		}
	} else {
		errorMsg := THUMBS_DOWN + "Date not set, use YYYY-MM-DD as date format\n"
		return CommandResult{Entry: InvalidEntry{Error: errorMsg}}
	}
}

func (whiteboard QuietWhiteboardApp) handleUpdateCommand(input string, context SlackContext) CommandResult {
	if err := handleEmptyInput(input); err != nil {
		return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
	}

	if entryType, ok := whiteboard.EntryMap[context.User.Username]; ok {
		entryType.GetEntry().Title = input

		if _, err := whiteboard.Repository.SaveEntry(entryType); err != nil {
			return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
		}

		return CommandResult{Entry: entryType}
	} else {
		return CommandResult{Entry: InvalidEntry{Error: MISSING_ENTRY}}
	}
}

func (whiteboard QuietWhiteboardApp) handleCreateCommand(input string, context SlackContext, factory EntryFactory) CommandResult {
	var entry fmt.Stringer

	if err := handleEmptyInput(input); err != nil {
		return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
	}

	standup, _ := whiteboard.Store.GetStandup(context.Channel.ChannelId)
	entryType := factory(whiteboard.Clock, context.User.Author, input, standup)
	whiteboard.EntryMap[context.User.Username] = entryType
	postResult, err := whiteboard.Repository.SaveEntry(entryType)
	if err != nil {
		return CommandResult{Entry: InvalidEntry{Error: err.Error()}}
	}
	entryType.GetEntry().Id = postResult.ItemId
	entry = *entryType.GetEntry()
	return CommandResult{Entry: entry}

}

func handleEmptyInput(input string) (err error) {
	if len(input) == 0 {
		err = errors.New(MISSING_INPUT)
	}
	return
}
