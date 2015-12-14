package app

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	. "github.com/xtreme-andleung/whiteboardbot/rest"
	"math/rand"
	"strings"
	"time"
)

const (
	botName        = "whiteboard-bot"
	usage   string = "*Usage*:\n" +
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
	input := ev.Text

	command, input := readNextCommand(input)
	if !matches(command, "wb") {
		return
	}

	username, author, ok := getSlackUser(slackClient, ev.User)
	if !ok {
		return
	}

	entryType := entryMap[username]

	command, input = readNextCommand(input)
	switch {
	case matches(command, "?"):
		postMarkdownMessageToSlack(usage, slackClient, ev.Channel)
		return
	case matches(command, "faces"):
		entryMap[username] = NewFace(clock, author, input)
	case matches(command, "interestings"):
		entryMap[username] = NewInteresting(clock, author, input)
	case matches(command, "helps"):
		entryMap[username] = NewHelp(clock, author, input)
	case matches(command, "events"):
		entryMap[username] = NewEvent(clock, author, input)
	case matches(command, "name") || matches(command, "title"):
		entryType.GetEntry().Title = input
	case matches(command, "body"):
		switch entryType.(type) {
		default:
			entryType.GetEntry().Body = input
		case Face:
			postMessageToSlack("Face does not have a body! " + randomInsult(), slackClient, ev.Channel)
			return
		}
	case matches(command, "date"):
		if parsedDate, err := time.Parse("2006-01-02", input); err == nil {
			entryType.GetEntry().Date = parsedDate
		} else {
			postMessageToSlack(entryType.String() + "\nDate not set, use YYYY-MM-DD as date format", slackClient, ev.Channel)
			return
		}
	default:
		postMessageToSlack(fmt.Sprintf("%v no you %v", username, ev.Text[3:]), slackClient, ev.Channel)
		return
	}

	entryType = entryMap[username]
	output := entryType.String()
	if entryType.Validate() {
		if itemId, ok := postEntryToWhiteboard(restClient, entryType); ok {
			output = appendStatus(entryType, output)
			entryType.GetEntry().Id = itemId
		}
	}
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

func postMessageToSlack(message string, slackClient SlackClient, channel string) {
	postToSlack(message, slackClient, channel, slack.PostMessageParameters{})
}

func postMarkdownMessageToSlack(message string, slackClient SlackClient, channel string) {
	params := slack.PostMessageParameters{Markdown: true}
	postToSlack(message, slackClient, channel, params)
}

func postToSlack(message string, slackClient SlackClient, channel string, params slack.PostMessageParameters) {
	fmt.Printf("Posting message to slack:\n%v\n", message)
	params.Username = botName
	slackClient.PostMessage(channel, message, params)
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
		fmt.Printf("%v, %v\n", username, err)
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

func appendStatus(entryType EntryType, output string) string {
	if isExistingEntry(entryType.GetEntry()) {
		return output + "\nitem updated"
	} else {
		return output + "\nitem created"
	}
}

func readNextCommand(input string) (keyword string, newInput string) {
	index := strings.Index(input, " ")
	if index == -1 {
		index = len(input)
	}
	keyword = strings.ToLower(input[:index])
	newInput = strings.TrimPrefix(input[index:], " ")
	return
}