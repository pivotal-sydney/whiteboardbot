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

func ParseMessageEvent(slackClient SlackClient, restClient RestClient, clock Clock, ev *slack.MessageEvent) {
	if !strings.HasPrefix(strings.ToLower(ev.Text), "wb ") {
		return
	}

	username, author, ok := getSlackUser(slackClient, ev.User)
	if !ok {
		return
	}

	input := ev.Text[3:]

	entryType, ok := entryMap[username]
	if !ok {
		entryMap[username] = &Entry{}
		entryType = entryMap[username]
	}

	index := strings.Index(input, " ")
	if index == -1 {
		index = len(input)
	}

	keyword := strings.ToLower(input[:index])
	switch {
	case matches(keyword, "?"):
		postMarkdownMessageToSlack(usage, slackClient, ev.Channel)
		return
	case matches(keyword, "faces"):
		entryType = NewFace(clock, author)
		entryMap[username] = entryType
		populateEntry(input, index, entryType)
	case matches(keyword, "interestings"):
		entryType = NewInteresting(clock, author)
		entryMap[username] = entryType
		populateEntry(input, index, entryType)
	case matches(keyword, "helps"):
		entryType = NewHelp(clock, author)
		entryMap[username] = entryType
		populateEntry(input, index, entryType)
	case matches(keyword, "events"):
		entryType = NewEvent(clock, author)
		entryMap[username] = entryType
		populateEntry(input, index, entryType)
	case matches(keyword, "name") || matches(keyword, "title"):
		entryType.GetEntry().Title = input[index + 1:]
	case matches(keyword, "body"):
		switch entryType.(type) {
		default:
			entryType.GetEntry().Body = input[index + 1:]
		case Face:
			postMessageToSlack("Face does not have a body! " + randomInsult(), slackClient, ev.Channel)
			return
	}
	case matches(keyword, "date"):
		if parsedDate, err := time.Parse("2006-01-02", input[index + 1:]); err == nil {
			entryType.GetEntry().Date = parsedDate
		} else {
			postMessageToSlack(entryType.String() + "\nDate not set, use YYYY-MM-DD as date format", slackClient, ev.Channel)
			return
		}
	default:
		postMessageToSlack(fmt.Sprintf("%v no you %v", username, input), slackClient, ev.Channel)
		return
	}

	output := entryType.String()
	if entryType.Validate() {
		itemId, ok := postEntryToWhiteboard(restClient, entryType)
		if ok {
			output = appendStatus(entryType, output)
			entryType.GetEntry().Id = itemId
		}
	}
	fmt.Printf("Posting message: %v\n", output)
	postMessageToSlack(output, slackClient, ev.Channel)
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

func getSlackUser(slackClient SlackClient, eventUser string) (username, author string, ok bool) {
	if slackUser, err := slackClient.GetUserInfo(eventUser); err == nil {
		username = slackUser.Name
		author = GetAuthor(slackUser)
		ok = true
	} else {
		fmt.Printf("%v, %v", username, err)
		ok = false
	}
	return
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