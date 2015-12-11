package app

import (
	"github.com/nlopes/slack"
	"strings"
	"fmt"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	. "github.com/xtreme-andleung/whiteboardbot/rest"
	"time"
	"math/rand"
)

const (
	botName = "whiteboard-bot"
	usage string = "*Usage*:\n" +
		"    `wb [command] [text...]`\n" +
		"where commands include:\n" +
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

var entryMap = make(map[string]EntryType)

var insults = [...]string{"Stupid.", "You idiot.", "You fool."}

func init() {
	rand.Seed(7483658374658473)
}

type SlackClient interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}

func ParseMessageEvent(slackClient SlackClient, restClient RestClient, clock Clock, ev *slack.MessageEvent) (username string, message string) {
	if strings.HasPrefix(strings.ToLower(ev.Text), "wb ") {
		user, err := slackClient.GetUserInfo(ev.User)
		if err != nil {
			fmt.Printf("%v, %v", ev.User, err)
			return
		}
		username = user.Name
		message = ev.Text[3:]

		entryType, ok := entryMap[username]
		if !ok {
			entryMap[username] = &Entry{}
			entryType = entryMap[username]
		}

		index := strings.Index(message, " ")
		if index == -1 {
			index = len(message)
		}

		keyword := strings.ToLower(message[:index])
		switch {
		case matches(keyword, "?"):
			message = usage
			postMarkdownMessageToSlack(usage, slackClient, ev.Channel)
			return
		case matches(keyword, "faces"):
			entryType = NewFace(clock, GetAuthor(user))
			entryMap[username] = entryType
			populateEntry(message, index, entryType)
		case matches(keyword, "interestings"):
			entryType = NewInteresting(clock, GetAuthor(user))
			entryMap[username] = entryType
			populateEntry(message, index, entryType)
		case matches(keyword, "helps"):
			entryType = NewHelp(clock, GetAuthor(user))
			entryMap[username] = entryType
			populateEntry(message, index, entryType)
		case matches(keyword, "events"):
			entryType = NewEvent(clock, GetAuthor(user))
			entryMap[username] = entryType
			populateEntry(message, index, entryType)
		case matches(keyword, "name") || matches(keyword, "title"):
			entryType.GetEntry().Title = message[index + 1:]
		case matches(keyword, "body"):
			switch entryType.(type) {
			default:
				entryType.GetEntry().Body = message[index + 1:]
			case Face:
				message = "Face does not have a body! " + randomInsult();
				postMessageToSlack(message, slackClient, ev.Channel)
				return
		}
		case matches(keyword, "date"):
			parsedDate, err := time.Parse("2006-01-02", message[index + 1:])
			if err != nil {
				message = entryType.String() + "\nDate not set, use YYYY-MM-DD as date format"
				postMessageToSlack(message, slackClient, ev.Channel)
				return
			} else {
				entryType.GetEntry().Date = parsedDate
			}
		default:
			message = fmt.Sprintf("%v no you %v", user.Name, message)
			postMessageToSlack(message, slackClient, ev.Channel)
			return
		}

		message = entryType.String()
		if entryType.Validate() {
			itemId, ok := postEntryToWhiteboard(restClient, entryType)
			if ok {
				message = appendStatus(entryType, message)
				entryType.GetEntry().Id = itemId
			}
		}
		fmt.Printf("Posting message: %v", message)
		postMessageToSlack(message, slackClient, ev.Channel)
	}
	return
}

func matches(keyword string, command string) bool {
	return len(keyword) <= len(command) && command[:len(keyword)] == keyword
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

func populateEntry(message string, index int, entryType EntryType) {
	entryType.GetEntry().Title = strings.TrimPrefix(message[index:], " ")
}

func postMessageToSlack(message string, slackClient SlackClient, channel string) {
	slackClient.PostMessage(channel, message, slack.PostMessageParameters{Username: botName})
}

func postMarkdownMessageToSlack(message string, slackClient SlackClient, channel string) {
	slackClient.PostMessage(channel, message, slack.PostMessageParameters{Markdown: true, Username: botName})
}

func postEntryToWhiteboard(restClient RestClient, entryType EntryType) (itemId string, ok bool) {
	var request = createRequest(entryType, isExistingEntry(entryType.GetEntry()))
	itemId, ok = restClient.Post(request)
	return
}

func randomInsult() string {
	return insults[rand.Intn(len(insults))]
}

func GetAuthor(user *slack.User) (realName string) {
	realName = user.Profile.RealName
	if len(realName) == 0 {
		realName = user.Name
	}
	return
}

func appendStatus(entryType EntryType, message string) string {
	if isExistingEntry(entryType.GetEntry()) {
		return message + "\nitem updated"
	} else {
		return message + "\nitem created"
	}
}