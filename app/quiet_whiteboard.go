package app

import (
	"errors"
	"fmt"
	. "github.com/pivotal-sydney/whiteboardbot/model"
)

type QuietWhiteboard interface {
	ProcessCommand(string, SlackContext) (CommandResult, error)
	PostEntry(EntryType) (PostResult, error)
}

type CommandHandler func(input string, context SlackContext) (CommandResult, error)

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
	var entry fmt.Stringer

	if entry = resultIfEmptyTitle(input); entry == nil {
		standup, _ := whiteboard.Store.GetStandup(context.Channel.ChannelId)
		face := NewFace(whiteboard.Clock, context.User.Author, input, standup)
		whiteboard.EntryMap[context.User.Username] = face
		if _, err := whiteboard.PostEntry(face); err != nil {
			return CommandResult{Entry: InvalidEntry{Error: err.Error()}}, nil
		}
		entry = *face.GetEntry()
	}

	return CommandResult{Entry: entry}, nil
}

func resultIfEmptyTitle(input string) fmt.Stringer {
	if len(input) == 0 {
		return InvalidEntry{Error: THUMBS_DOWN + "Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? Try `wb ?`"}
	}
	return nil
}
