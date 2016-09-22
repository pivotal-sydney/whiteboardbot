package app

import (
	"errors"
	"fmt"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"time"
)

type QuietWhiteboard interface {
	ProcessCommand(string, SlackContext) (CommandResult, error)
	PostEntry(EntryType) (PostResult, error)
}

type CommandHandler func(input string, context SlackContext) (CommandResult, error)

type EntryFactory func(clock Clock, author, title string, standup Standup) EntryType

type QuietWhiteboardApp struct {
	Clock      Clock
	RestClient RestClient
	Store      Store
	CommandMap map[string]CommandHandler
	EntryMap   map[string]EntryType
}

type CommandResult struct {
	Entry fmt.Stringer
}

type PostResult struct {
	ItemId string
}

func NewQuietWhiteboard(restClient RestClient, store Store, clock Clock) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{
		Clock:      clock,
		RestClient: restClient,
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
}

func (whiteboard QuietWhiteboardApp) ProcessCommand(input string, context SlackContext) (CommandResult, error) {
	command, input := readNextCommand(input)
	return whiteboard.handleCommand(command, input, context)
}

func (whiteboard QuietWhiteboardApp) PostEntry(entryType EntryType) (PostResult, error) {
	itemId, ok := PostEntryToWhiteboard(whiteboard.RestClient, entryType)

	if !ok {
		return PostResult{}, errors.New("Problem creating post.")
	}

	return PostResult{itemId}, nil
}

func (whiteboard QuietWhiteboardApp) handleCommand(command, input string, context SlackContext) (CommandResult, error) {
	for key := range whiteboard.CommandMap {
		if matches(command, key) {
			callback := whiteboard.CommandMap[key]
			return callback(input, context)
		}
	}

	return CommandResult{Entry: InvalidEntry{Error: "Ooops"}}, nil
}

func (whiteboard QuietWhiteboardApp) registerCommand(command string, callback CommandHandler) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string, _ SlackContext) (CommandResult, error) {
	return CommandResult{Entry: TextEntry{Text: USAGE}}, nil
}

func (whiteboard QuietWhiteboardApp) handleRegistrationCommand(standupId string, context SlackContext) (command CommandResult, err error) {
	command = CommandResult{}

	standup, ok := whiteboard.RestClient.GetStandup(standupId)
	if !ok {
		err = errors.New("Standup not found!")
		return
	}

	whiteboard.Store.SetStandup(context.Channel.ChannelId, standup)

	text := fmt.Sprintf("Standup %v has been registered! You can now start creating Whiteboard entries!", standup.Title)

	command.Entry = TextEntry{Text: text}
	return
}

func (whiteboard QuietWhiteboardApp) handleFacesCommand(input string, context SlackContext) (CommandResult, error) {
	return whiteboard.handleCreateCommand(input, context, NewFace)
}

func (whiteboard QuietWhiteboardApp) handleHelpsCommand(input string, context SlackContext) (CommandResult, error) {
	return whiteboard.handleCreateCommand(input, context, NewHelp)
}

func (whiteboard QuietWhiteboardApp) handleInterestingsCommand(input string, context SlackContext) (CommandResult, error) {
	return whiteboard.handleCreateCommand(input, context, NewInteresting)
}

func (whiteboard QuietWhiteboardApp) handleEventsCommand(input string, context SlackContext) (CommandResult, error) {
	return whiteboard.handleCreateCommand(input, context, NewEvent)
}

func (whiteboard QuietWhiteboardApp) handleBodyCommand(input string, context SlackContext) (CommandResult, error) {

	if len(input) == 0 {
		errorMsg := THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb b My title`\nNeed help? Try `wb ?`"
		return CommandResult{Entry: InvalidEntry{Error: errorMsg}}, nil
	}

	username := context.User.Username

	if entryType, ok := whiteboard.EntryMap[username]; ok {
		entry := *entryType.GetEntry()

		if entry.ItemKind == "New face" {
			errorMsg := ":-1:\nHey, new faces should not have a body!"
			return CommandResult{InvalidEntry{Error: errorMsg}}, nil
		}

		entry.Body = input
		whiteboard.EntryMap[username] = entry

		entryType = whiteboard.EntryMap[username]

		if _, err := whiteboard.PostEntry(entryType); err != nil {
			return CommandResult{Entry: InvalidEntry{Error: err.Error()}}, nil
		}
		return CommandResult{Entry: entry}, nil
	} else {
		errorMsg := THUMBS_DOWN + "Hey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!"

		return CommandResult{Entry: InvalidEntry{Error: errorMsg}}, nil
	}
}

func (whiteboard QuietWhiteboardApp) handleDateCommand(input string, context SlackContext) (CommandResult, error) {

	if len(input) == 0 {
		errorMsg := THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb d 2017-05-21`\nNeed help? Try `wb ?`"
		return CommandResult{Entry: InvalidEntry{Error: errorMsg}}, nil
	}

	if parsedDate, err := time.Parse(DATE_FORMAT, input); err == nil {
		if entryType, ok := whiteboard.EntryMap[context.User.Username]; ok {
			entryType.GetEntry().Date = parsedDate.Format(DATE_FORMAT)

			request := createRequest(entryType, len(entryType.GetEntry().Id) > 0)

			whiteboard.RestClient.Post(request)

			return CommandResult{Entry: entryType.GetEntry()}, nil
		} else {
			errorMsg := THUMBS_DOWN + "Hey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!"

			return CommandResult{Entry: InvalidEntry{Error: errorMsg}}, nil
		}
	} else {
		errorMsg := THUMBS_DOWN + "Date not set, use YYYY-MM-DD as date format\n"
		return CommandResult{Entry: InvalidEntry{Error: errorMsg}}, nil
	}
}

func (whiteboard QuietWhiteboardApp) handleCreateCommand(input string, context SlackContext, factory EntryFactory) (CommandResult, error) {
	var entry fmt.Stringer

	if entry = resultIfEmptyTitle(input); entry == nil {
		standup, _ := whiteboard.Store.GetStandup(context.Channel.ChannelId)
		entryType := factory(whiteboard.Clock, context.User.Author, input, standup)
		whiteboard.EntryMap[context.User.Username] = entryType
		postResult, err := whiteboard.PostEntry(entryType)
		if err != nil {
			return CommandResult{Entry: InvalidEntry{Error: err.Error()}}, nil
		}
		entryType.GetEntry().Id = postResult.ItemId
		entry = *entryType.GetEntry()
	}

	return CommandResult{Entry: entry}, nil
}

func resultIfEmptyTitle(input string) fmt.Stringer {
	if len(input) == 0 {
		return InvalidEntry{Error: THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"}
	}
	return nil
}
