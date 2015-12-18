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
)

const (
	botName        = "whiteboardbot"
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

var entryMap = make(map[string]EntryType)

type WhiteboardApp struct {
	SlackClient SlackClient
	RestClient  RestClient
	Clock       Clock
	Store       Store
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
			whiteboard.Store.Set(ev.Channel, int(standupId))
			postMessageToSlack(fmt.Sprintf("Standup Id: %v has been registered! You can now start creating Whiteboard entries!", standupId), whiteboard.SlackClient, ev.Channel)
		} else {
			handleRegisterationFailure(whiteboard.SlackClient, ev.Channel)
		}
		return
	}

	standupId, ok := whiteboard.Store.Get(ev.Channel)
	if !ok {
		handleRegisterationFailure(whiteboard.SlackClient, ev.Channel)
		return
	}

	username, author, ok := getSlackUser(whiteboard.SlackClient, ev.User)
	if !ok {
		return
	}

	entryType := entryMap[username]
	switch {
	case matches(command, "?"):
		postMarkdownMessageToSlack(usage, whiteboard.SlackClient, ev.Channel)
		return
	case matches(command, "faces"):
		entryMap[username] = NewFace(whiteboard.Clock, author, input, standupId)
	case matches(command, "interestings"):
		entryMap[username] = NewInteresting(whiteboard.Clock, author, input, standupId)
	case matches(command, "helps"):
		entryMap[username] = NewHelp(whiteboard.Clock, author, input, standupId)
	case matches(command, "events"):
		entryMap[username] = NewEvent(whiteboard.Clock, author, input, standupId)
	case matches(command, "name") || matches(command, "title"):
		entryType.GetEntry().Title = input
	case matches(command, "body"):
		switch entryType.(type) {
		default:
			entryType.GetEntry().Body = input
		case Face:
			postMessageToSlack("Face does not have a body! "+randomInsult(), whiteboard.SlackClient, ev.Channel)
			return
		}
	case matches(command, "date"):
		if parsedDate, err := time.Parse("2006-01-02", input); err == nil {
			entryType.GetEntry().Date = parsedDate.Format("2006-01-02")
		} else {
			postMessageToSlack(entryType.String()+"\nDate not set, use YYYY-MM-DD as date format", whiteboard.SlackClient, ev.Channel)
			return
		}
	case matches(command, "present"):
		items, ok := whiteboard.RestClient.GetStandupItems(standupId)
		if ok {
			if items.Empty() {
				postMessageToSlack("Hey, there's no entries in today's standup yet, why not add some?", whiteboard.SlackClient, ev.Channel)
				return
			}
			postMessageToSlack(items.String(), whiteboard.SlackClient, ev.Channel)
			return
		}
	default:
		var userInput string
		if fileUpload {
			_, userInput = readNextCommand(ev.File.Title[2:])
		} else {
			_, userInput = readNextCommand(ev.Text[2:])
		}
		postMessageToSlack(fmt.Sprintf("%v no you %v", username, userInput), whiteboard.SlackClient, ev.Channel)
		return
	}

	entryType = entryMap[username]

	if fileUpload {
		entryType.GetEntry().Body = fmt.Sprintf("%v\n<img src=\"%v\" style=\"max-width: 500px\">", ev.File.InitialComment.Comment, ev.File.URL)
	}

	output := entryType.String()
	if entryType.Validate() {
		if itemId, ok := postEntryToWhiteboard(whiteboard.RestClient, entryType, standupId); ok {
			output = appendStatus(entryType, output)
			entryType.GetEntry().Id = itemId
		}
	}
	postMessageToSlack(output, whiteboard.SlackClient, ev.Channel)
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

func postEntryToWhiteboard(restClient RestClient, entryType EntryType, standupId int) (itemId string, ok bool) {
	var request = createRequest(entryType, isExistingEntry(entryType.GetEntry()))
	itemId, ok = restClient.Post(request, standupId)
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

func handleRegisterationFailure(slackClient SlackClient, channel string) {
	postMessageToSlack("You haven't registered your standup yet. wb register <id> first!  (or short wb r <id>)", slackClient, channel)
	return
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
