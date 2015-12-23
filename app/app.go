package app

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	"math/rand"
	"strings"
	"time"
	"regexp"
	"encoding/json"
)

const (
	Usage   string =
	"*Usage*:\n" +
	"        `wb [command] [text...]`\n" +
	"    where commands include:\n" +
	"*Registration Command*\n" +
	"        `register`, `r` - followed by <standup_id>, registers current channel to Whiteboard's standup id\n" +
	"\n" +
	"*Presentation Command*\n" +
	"		 `present`, `p` - presents today's standup\n" +
	"\n" +
	"*Create Commands*\n" +
	"        `faces`, `f` - followed by a title, creates a new faces entry\n" +
	"        `interestings`, `i` - followed by a title, creates a new interestings entry\n" +
	"        `helps`, `h` - followed by a title, creates a new helps entry\n" +
	"        `events`, `e` - followed by a title, creates a new events entry\n" +
	"\n" +
	"*Detail Commands* (updates details of a started entry)\n" +
	"        `title`, `t`, `name`, `n` - updates a name/title detail to a started entry\n" +
	"        `body`, `b` - updates a body detail to a started entry\n" +
	"        `date`, `d` - updates a date detail to a started entry (YYYY-MM-DD)\n" +
	"\n" +
	"Example:\n" +
	"        `wb f New Face!` - will create a new face with the name 'New Face!'\n" +
	"        `wb d 2015-01-02` - will update the new face date to 02 Jan 2015"
)

var insults = [...]string{"Stupid.", "You idiot.", "You fool."}

func init() {
	rand.Seed(7483658374658473)
}

func (whiteboard WhiteboardApp) ParseMessageEvent(ev *slack.MessageEvent) {
	input := ev.Text
	fileUpload := ev.Upload
	if fileUpload {
		input = ev.File.Title
	}

	command, input := readNextCommand(input)
	if !matches(command, "wb") {
		return
	}

	command, input = readNextCommand(input)

	if matches(command, "register") {
		whiteboard.registerStandup(input, ev.Channel)
		return
	}

	standupJson, ok := whiteboard.Store.Get(ev.Channel)
	if !ok {
		handleNotRegistered(whiteboard.SlackClient, ev.Channel)
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
		whiteboard.SlackClient.PostMessageWithMarkdown(Usage, ev.Channel, "")
		return
	case matches(command, "faces"):
		if len(input) == 0 {
			handleMissingTitle(whiteboard, ev.Channel)
			return
		}
		whiteboard.EntryMap[username] = NewFace(whiteboard.Clock, author, input, standup)
	case matches(command, "interestings"):
		if len(input) == 0 {
			handleMissingTitle(whiteboard, ev.Channel)
			return
		}
		whiteboard.EntryMap[username] = NewInteresting(whiteboard.Clock, author, input, standup)
	case matches(command, "helps"):
		if len(input) == 0 {
			handleMissingTitle(whiteboard, ev.Channel)
			return
		}
		whiteboard.EntryMap[username] = NewHelp(whiteboard.Clock, author, input, standup)
	case matches(command, "events"):
		if len(input) == 0 {
			handleMissingTitle(whiteboard, ev.Channel)
			return
		}
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
			whiteboard.SlackClient.PostMessage("Face does not have a body! " + randomInsult(), ev.Channel, THUMBS_DOWN)
			return
		}
	case matches(command, "date"):
		if missingEntry(entryType) {
			handleMissingEntry(whiteboard.SlackClient, ev.Channel)
			return
		}

		if parsedDate, err := time.Parse(DATE_FORMAT, input); err == nil {
			entryType.GetEntry().Date = parsedDate.Format(DATE_FORMAT)
		} else {
			whiteboard.SlackClient.PostEntry(entryType, ev.Channel, "Date not set, use YYYY-MM-DD as date format\n")
			return
		}
	case matches(command, "present"):
		items, ok := whiteboard.RestClient.GetStandupItems(standup.Id)
		if ok {
			if items.Empty() {
				whiteboard.SlackClient.PostMessage("Hey, there's no entries in today's standup yet, why not add some?", ev.Channel, THUMBS_DOWN)
				return
			}
			whiteboard.SlackClient.PostMessage(items.String(), ev.Channel, "")
			return
		}
	default:
		var userInput string
		if fileUpload {
			_, userInput = readNextCommand(ev.File.Title[2:])
		} else {
			_, userInput = readNextCommand(ev.Text[2:])
		}
		whiteboard.SlackClient.PostMessage(fmt.Sprintf("%v no you %v", username, userInput), ev.Channel, "")
		return
	}

	entryType = whiteboard.EntryMap[username]

	if fileUpload {
		entryType.GetEntry().Body = fmt.Sprintf("%v\n<img src=\"%v\" style=\"max-width: 500px\">", ev.File.InitialComment.Comment, ev.File.URL)
	}

	status := ""
	if entryType.Validate() {
		if itemId, ok := PostEntryToWhiteboard(whiteboard.RestClient, entryType, standup.Id); ok {
			status = THUMBS_UP
			entryType.GetEntry().Id = itemId
		}
	}
	whiteboard.SlackClient.PostEntry(entryType, ev.Channel, status)
	return
}

func matches(keyword string, command string) bool {
	return len(keyword) > 0 && len(keyword) <= len(command) && command[:len(keyword)] == keyword
}

func randomInsult() string {
	return insults[rand.Intn(len(insults))]
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

func missingEntry(entryType EntryType) bool {
	return entryType == nil
}
