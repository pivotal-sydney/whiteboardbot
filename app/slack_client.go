package app

import (
	"github.com/nlopes/slack"
	"strings"
	"fmt"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	. "github.com/xtreme-andleung/whiteboardbot/rest"
	"time"
)

var entryType EntryType
var entry *Entry

type SlackClient interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}

func ParseMessageEvent(slackClient SlackClient, restClient RestClient, clock Clock, ev *slack.MessageEvent) (username string, message string) {
	if strings.HasPrefix(ev.Text, "wb ") {
		user, err := slackClient.GetUserInfo(ev.User)
		if err != nil {
			fmt.Printf("%v, %v", ev.User, err)
			return
		}
		username = user.Name
		message = ev.Text[3:]
		if entryType != nil {
			switch entryType.(type) {
			case Face:
				entry = entryType.(Face).Entry
			case Interesting:
				entry = entryType.(Interesting).Entry
			case Event:
				entry = entryType.(Event).Entry
			case Help:
				entry = entryType.(Help).Entry
			}
		}
		if strings.HasPrefix(message, "faces") {
			entryType = NewFace(clock, username)
		} else if strings.HasPrefix(message, "interestings") {
			entryType = NewInteresting(clock, username)
		} else if strings.HasPrefix(message, "events") {
			entryType = NewEvent(clock, username)
		} else if strings.HasPrefix(message, "helps") {
			entryType = NewHelp(clock, username)
		} else if strings.HasPrefix(message, "name ") {
			entry.Title = message[5:]
		} else if strings.HasPrefix(message, "title ") {
			entry.Title = message[6:]
		} else if strings.HasPrefix(message, "body ") {
			entry.Body = message[5:]
		} else if strings.HasPrefix(message, "date ") {
			parsedDate, err := time.Parse("2006-01-02", message[5:])
			if err != nil {
				message = entryType.String() + "\nDate not set, use YYYY-MM-DD as date format"
				slackClient.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
				return
			} else {
				entry.Time = parsedDate
			}
		} else {
			message = fmt.Sprintf("%v no you %v", user.Name, message)
			slackClient.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
			return
		}
		message = entryType.String()
		if entryType.Validate() {
			var request = createRequest(entryType, entry)
			itemId, ok := restClient.Post(request)
			if ok {
				entry.Id = itemId
				if len(request.Id) > 0 {
					message += "\nitem updated"
				} else {
					message += "\nitem created"
				}
			}
		}
		fmt.Printf("Posting message: %v", message)
		slackClient.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
	}
	return
}

func createRequest(entryType EntryType, entry *Entry) (request WhiteboardRequest) {
	if len(entry.Id) > 0 {
		request = entryType.MakeUpdateRequest()
	} else {
		request = entryType.MakeCreateRequest()
	}
	return
}