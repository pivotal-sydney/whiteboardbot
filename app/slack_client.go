package app

import (
	"github.com/nlopes/slack"
	"strings"
	"fmt"
	. "github.com/xtreme-andleung/whiteboardbot/entry"
	"time"
)

var face Face

type SlackClient interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
	GetUserInfo(user string) (*slack.User, error)
}

func ParseMessageEvent(client SlackClient, clock Clock, ev *slack.MessageEvent) (username string, message string) {
	if strings.HasPrefix(ev.Text, "wb ") {
		user, err := client.GetUserInfo(ev.User)
		if err != nil {
			fmt.Printf("%v, %v", ev.User, err)
			return
		}
		username = user.Name
		message = ev.Text[3:]
		if strings.HasPrefix(message, "faces") {
			face = NewFace(clock)
			message = face.String()
		} else if strings.HasPrefix(message, "name ") {
			face.Name = message[5:]
			message = face.String()
		} else if strings.HasPrefix(message, "date ") {
			parsedDate, err := time.Parse("2006-01-02", message[5:])
			if err != nil {
				message = face.String() + "\nDate not set, use YYYY-MM-DD as date format"
			} else {
				face.Time = parsedDate
				message = face.String()
			}
		} else {
			message = strings.Join([]string{user.Name, "no you", message}, " ")
		}
		fmt.Printf("Posting message: %v", message)
		client.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
	}
	return
}