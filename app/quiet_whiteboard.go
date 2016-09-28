package app

import (
	"errors"
	"fmt"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"strconv"
	"strings"
	"time"
)

type QuietWhiteboard interface {
	ProcessCommand(string, SlackContext) CommandResult
}

type CommandHandler func(input string, context SlackContext) CommandResult

type EntryFactory func(clock Clock, author, title string, standup Standup) EntryType

type QuietWhiteboardApp struct {
	Clock      Clock
	Repository StandupRepository
	Store      Store
	CommandMap map[string]CommandHandler
	EntryMap   map[string]EntryType
}

type CommandResult fmt.Stringer

type EntryCommandResult struct {
	Title    string
	Status   string
	HelpText string
	Entry    fmt.Stringer
}

type MessageCommandResult struct {
	Text   string
	Status string
}

func (r MessageCommandResult) String() string {
	status := ""
	if r.Status != "" {
		status = r.Status
	}

	return fmt.Sprintf("%s%s", status, r.Text)
}

func (r EntryCommandResult) String() string {
	helpText := ""
	if r.HelpText != "" {
		helpText = r.HelpText + "\n"
	}

	return fmt.Sprintf("%s%s\n%s%s", r.Status, r.Title, helpText, r.Entry.String())
}

func NewQuietWhiteboard(gateway StandupRepository, store Store, clock Clock) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{
		Clock:      clock,
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
	whiteboard.registerCommand("present", whiteboard.handlePresentCommand)
}

func (whiteboard QuietWhiteboardApp) ProcessCommand(input string, context SlackContext) CommandResult {
	command, input := ReadNextCommand(input)
	return whiteboard.handleCommand(command, input, context)
}

func (whiteboard QuietWhiteboardApp) handleCommand(command, input string, context SlackContext) CommandResult {
	for key := range whiteboard.CommandMap {
		if matches(command, key) {
			callback := whiteboard.CommandMap[key]
			return callback(input, context)
		}
	}

	return MessageCommandResult{Text: "Ooops"}
}

func (whiteboard QuietWhiteboardApp) registerCommand(command string, callback CommandHandler) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string, _ SlackContext) CommandResult {
	return MessageCommandResult{Text: USAGE}
}

func (whiteboard QuietWhiteboardApp) handleRegistrationCommand(standupId string, context SlackContext) CommandResult {
	standup, err := whiteboard.Repository.FindStandup(standupId)
	if err != nil {
		return MessageCommandResult{Status: THUMBS_DOWN, Text: "Standup not found!"}
	}

	whiteboard.Store.SetStandup(context.Channel.Id, standup)

	text := fmt.Sprintf("Standup %v has been registered! You can now start creating Whiteboard entries!", standup.Title)

	return MessageCommandResult{Status: THUMBS_UP, Text: text}
}

func (whiteboard QuietWhiteboardApp) handlePresentCommand(numDays string, context SlackContext) CommandResult {
	standup, ok := whiteboard.Store.GetStandup(context.Channel.Id)
	if !ok {
		return MessageCommandResult{Status: THUMBS_DOWN, Text: MISSING_STANDUP}
	}

	standupId := strconv.Itoa(standup.Id)

	standupItems, err := whiteboard.Repository.GetStandupItems(standupId)
	if err != nil {
		return MessageCommandResult{Text: err.Error(), Status: THUMBS_DOWN}
	}

	numberOfDays, _ := strconv.Atoi(numDays)
	if numberOfDays > 0 {
		standupItems = standupItems.Filter(numberOfDays, whiteboard.Clock, context.User.TimeZone)
	}

	return MessageCommandResult{Text: standupItems.String()}
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
		return MessageCommandResult{Text: err.Error(), Status: THUMBS_DOWN}
	}

	username := context.User.Username

	if entryType, ok := whiteboard.EntryMap[username]; ok {
		entry := entryType.GetEntry()

		if entry.ItemKind == "New face" {
			errorMsg := "Hey, new faces should not have a body!"
			return MessageCommandResult{Text: errorMsg, Status: THUMBS_DOWN}
		}

		entryType.GetEntry().Body = input

		if _, err := whiteboard.Repository.SaveEntry(entryType); err != nil {
			return MessageCommandResult{Text: err.Error()}
		}
		return makeEntryCommandResult(entryType, false)
	} else {
		return MessageCommandResult{Text: MISSING_ENTRY, Status: THUMBS_DOWN}
	}
}

func (whiteboard QuietWhiteboardApp) handleDateCommand(input string, context SlackContext) CommandResult {
	if err := handleEmptyInput(input); err != nil {
		return MessageCommandResult{Text: err.Error(), Status: THUMBS_DOWN}
	}

	if parsedDate, err := time.Parse(DATE_FORMAT, input); err == nil {
		if entryType, ok := whiteboard.EntryMap[context.User.Username]; ok {
			entryType.GetEntry().Date = parsedDate.Format(DATE_FORMAT)

			if _, err := whiteboard.Repository.SaveEntry(entryType); err != nil {
				return MessageCommandResult{Text: err.Error()}
			}

			return EntryCommandResult{Entry: entryType.GetEntry()}
		} else {
			return MessageCommandResult{Text: MISSING_ENTRY, Status: THUMBS_DOWN}
		}
	} else {
		errorMsg := "Date not set, use YYYY-MM-DD as date format\n"
		return MessageCommandResult{Text: errorMsg, Status: THUMBS_DOWN}
	}
}

func (whiteboard QuietWhiteboardApp) handleUpdateCommand(input string, context SlackContext) CommandResult {
	if err := handleEmptyInput(input); err != nil {
		return MessageCommandResult{Text: err.Error(), Status: THUMBS_DOWN}
	}

	if entryType, ok := whiteboard.EntryMap[context.User.Username]; ok {
		entryType.GetEntry().Title = input

		if _, err := whiteboard.Repository.SaveEntry(entryType); err != nil {
			return MessageCommandResult{Text: err.Error()}
		}

		return makeEntryCommandResult(entryType, false)
	} else {
		return MessageCommandResult{Text: MISSING_ENTRY, Status: THUMBS_DOWN}
	}
}

func (whiteboard QuietWhiteboardApp) handleCreateCommand(input string, context SlackContext, factory EntryFactory) CommandResult {

	if err := handleEmptyInput(input); err != nil {
		return MessageCommandResult{Text: err.Error(), Status: THUMBS_DOWN}
	}

	standup, ok := whiteboard.Store.GetStandup(context.Channel.Id)
	if !ok {
		return MessageCommandResult{Text: MISSING_STANDUP, Status: THUMBS_DOWN}
	}

	entryType := factory(whiteboard.Clock, context.User.Author, input, standup)
	whiteboard.EntryMap[context.User.Username] = entryType
	postResult, err := whiteboard.Repository.SaveEntry(entryType)
	if err != nil {
		return MessageCommandResult{Text: err.Error()}
	}
	entryType.GetEntry().Id = postResult.ItemId

	return makeEntryCommandResult(entryType, true)
}

func makeEntryCommandResult(entryType EntryType, newEntry bool) EntryCommandResult {
	entry := entryType.GetEntry()
	itemKind := entry.ItemKind
	helpText := ""
	if itemKind != "New face" && newEntry {
		helpText = NEW_ENTRY_HELP_TEXT
	}

	return EntryCommandResult{
		Title:    strings.ToUpper(itemKind),
		Status:   THUMBS_UP,
		HelpText: helpText,
		Entry:    entry,
	}
}

func handleEmptyInput(input string) (err error) {
	if len(input) == 0 {
		err = errors.New(MISSING_INPUT)
	}
	return
}
