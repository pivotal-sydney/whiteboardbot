package app

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	. "github.com/xtreme-andleung/whiteboardbot/persistance"
	. "github.com/xtreme-andleung/whiteboardbot/rest"
	. "github.com/xtreme-andleung/whiteboardbot/slack_client"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"regexp"
	"encoding/json"
)

const (
	usage   string = "*Usage*:\n" +
		"    `wb [command] [text...]`\n" +
		"where commands include:\n" +
		"    *Registration Command*\n" +
		"        `register`, `r` - followed by <standup id> - registers current channel to whiteboard's standup ID\n" +
		"\n" +
	    "    *Presentation Command*\n" +
		"		 `present`, `p` - presents today's standup\n" +
		"\n" +
		"    *Create Commands*\n" +
		"        `faces`, `f` - creates a new faces entry\n" +
		"        `interestings`, `i` - creates a new interestings entry\n" +
		"        `helps`, `h` - creates a new helps entry\n" +
		"        `events`, `e` - creates a new events entry\n" +
		"\n" +
		"    *Detail Commands* (after creating an entry)\n" +
		"        `title`, `t`, `name`, `n` - adds a name/title detail to a started entry\n" +
		"        `body`, `b` - adds a body detail to a started entry\n" +
		"        `date`, `d` - adds a date detail to a started entry (YYYY-MM-DD)\n" +
		"\n" +
		"When adding text along with a Create Command, the title of the entry will be set to the text\n" +
		"Ex:\n" +
		"    `wb f New Face!` - will create a new face with the name 'New Face!'"
)

type WhiteboardApp struct {
	SlackClient SlackClient
	RestClient  RestClient
	Clock       Clock
	Store       Store
	EntryMap    map[string]EntryType
}

var insults = [...]string{"Stupid.", "You idiot.", "You fool."}

func init() {
	rand.Seed(7483658374658473)
}

func (whiteboard WhiteboardApp) ParseMessageEvent(ev *slack.MessageEvent) {
	input := ev.Text

	var fileUpload bool
	if ev.Upload {
		input = ev.File.Title
		fileUpload = true
	}

	command, input := readNextCommand(input)
	if !matches(command, "wb") {
		return
	}

	command, input = readNextCommand(input)

	if matches(command, "register") {
		standupId, err := strconv.ParseInt(input, 10, 64)
		if err == nil {
			whiteboard.registerStandup(int(standupId), ev.Channel)
		} else {
			handleRegisterationFailure(whiteboard.SlackClient, ev.Channel)
		}
		return
	}

	standupJson, ok := whiteboard.Store.Get(ev.Channel)
	if !ok {
		handleRegisterationFailure(whiteboard.SlackClient, ev.Channel)
		return
	}
	var standup Standup
	json.Unmarshal([]byte(standupJson), &standup)

	username, author, ok := whiteboard.SlackClient.GetUserDetails(ev.User)
	if !ok {
		return
	}

	entryType := whiteboard.EntryMap[username]
	switch {
	case matches(command, "?"):
		whiteboard.SlackClient.PostMessageWithMarkdown(usage, ev.Channel)
		return
	case matches(command, "faces"):
		whiteboard.EntryMap[username] = NewFace(whiteboard.Clock, author, input, standup)
	case matches(command, "interestings"):
		whiteboard.EntryMap[username] = NewInteresting(whiteboard.Clock, author, input, standup)
	case matches(command, "helps"):
		whiteboard.EntryMap[username] = NewHelp(whiteboard.Clock, author, input, standup)
	case matches(command, "events"):
		whiteboard.EntryMap[username] = NewEvent(whiteboard.Clock, author, input, standup)
	case matches(command, "name") || matches(command, "title"):
		if missingEntry(entryType) {
			handleMissingEntry(whiteboard.SlackClient, ev.Channel)
			return
		}
		entryType.GetEntry().Title = input
	case matches(command, "body"):
		if missingEntry(entryType) {
			handleMissingEntry(whiteboard.SlackClient, ev.Channel)
			return
		}
		switch entryType.(type) {
		default:
			entryType.GetEntry().Body = input
		case Face:
			whiteboard.SlackClient.PostMessage("Face does not have a body! "+randomInsult(), ev.Channel)
			return
		}
	case matches(command, "date"):
		if missingEntry(entryType) {
			handleMissingEntry(whiteboard.SlackClient, ev.Channel)
			return
		}

		if parsedDate, err := time.Parse("2006-01-02", input); err == nil {
			entryType.GetEntry().Date = parsedDate.Format("2006-01-02")
		} else {
			whiteboard.SlackClient.PostMessage(entryType.String()+"\nDate not set, use YYYY-MM-DD as date format", ev.Channel)
			return
		}
	case matches(command, "present"):
		items, ok := whiteboard.RestClient.GetStandupItems(standup.Id)
		if ok {
			if items.Empty() {
				whiteboard.SlackClient.PostMessage("Hey, there's no entries in today's standup yet, why not add some?", ev.Channel)
				return
			}
			whiteboard.SlackClient.PostMessage(items.String(), ev.Channel)
			return
		}
	default:
		var userInput string
		if fileUpload {
			_, userInput = readNextCommand(ev.File.Title[2:])
		} else {
			_, userInput = readNextCommand(ev.Text[2:])
		}
		whiteboard.SlackClient.PostMessage(fmt.Sprintf("%v no you %v", username, userInput), ev.Channel)
		return
	}

	entryType = whiteboard.EntryMap[username]

	if fileUpload {
		entryType.GetEntry().Body = fmt.Sprintf("%v\n<img src=\"%v\" style=\"max-width: 500px\">", ev.File.InitialComment.Comment, ev.File.URL)
	}

	status := ""
	if entryType.Validate() {
		if itemId, ok := postEntryToWhiteboard(whiteboard.RestClient, entryType, standup.Id); ok {
			status = entryStatus(entryType)
			entryType.GetEntry().Id = itemId
		}
	}
	whiteboard.SlackClient.PostEntry(entryType, ev.Channel, status)
	return
}

func matches(keyword string, command string) bool {
	return len(keyword) > 0 && len(keyword) <= len(command) && command[:len(keyword)] == keyword
}

func isExistingEntry(entry *Entry) bool {
	return entry != nil && len(entry.Id) > 0
}

func createRequest(entryType EntryType, existingEntry bool) (request WhiteboardRequest) {
	if existingEntry {
		request = entryType.MakeUpdateRequest()
	} else {
		request = entryType.MakeCreateRequest()
	}
	return
}

func postEntryToWhiteboard(restClient RestClient, entryType EntryType, standupId int) (itemId string, ok bool) {
	var request = createRequest(entryType, isExistingEntry(entryType.GetEntry()))
	itemId, ok = restClient.Post(request, standupId)
	return
}

func randomInsult() string {
	return insults[rand.Intn(len(insults))]
}

func entryStatus(entryType EntryType) string {
	if isExistingEntry(entryType.GetEntry()) {
		return "\nitem updated"
	} else {
		return "\nitem created"
	}
}

func readNextCommand(input string) (keyword string, newInput string) {
	re := regexp.MustCompile("\\s+")
	loc := re.FindStringIndex(input)
	if loc != nil {
		keyword = strings.ToLower(input[:loc[0]])
		newInput = input[loc[1]:]
	} else {
		keyword = strings.ToLower(input)
		newInput = ""
	}
	return
}

func handleMissingEntry(slackClient SlackClient, channel string) {
	slackClient.PostMessageWithMarkdown("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event]` first!", channel)
}

func handleRegisterationFailure(slackClient SlackClient, channel string) {
	slackClient.PostMessage("You haven't registered your standup yet. wb register <id> first!  (or short wb r <id>)", channel)
	return
}

func missingEntry(entryType EntryType) bool {
	return entryType == nil
}
func (whiteboard WhiteboardApp) registerStandup(standupId int, channel string) {
	standup, ok := whiteboard.RestClient.GetStandup(standupId)
	if !ok {
		handleRegisterationFailure(whiteboard.SlackClient, channel)
		return
	}
	standupJson, _ := json.Marshal(standup)
	whiteboard.Store.Set(channel, string(standupJson))
	whiteboard.SlackClient.PostMessage(fmt.Sprintf("Standup Id: %v has been registered! You can now start creating Whiteboard entries!", standupId), channel)
}