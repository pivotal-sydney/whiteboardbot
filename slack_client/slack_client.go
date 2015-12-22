package slack_client
import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/model"
)

const botName = "whiteboardbot"

type Slack struct {
	SlackWrapper SlackWrapper
}

type SlackClient interface {
	PostMessage(message string, channel string)
	PostMessageWithMarkdown(message string, channel string)
	PostEntry(entryType model.EntryType, channel string, status string)
	GetUserDetails(user string) (username, author string, ok bool)
}

func (slackClient *Slack) PostMessage(message string, channel string) {
	slackClient.postMessage(message, channel, slack.PostMessageParameters{})
}

func (slackClient *Slack) PostMessageWithMarkdown(message string, channel string) {
	slackClient.postMessage(message, channel, slack.PostMessageParameters{Markdown: true})
}

func (slackClient *Slack) PostEntry(entryType model.EntryType, channel string, status string) {
	message := entryType.String() + status
	slackClient.PostMessage(message, channel)
}

func (slackClient *Slack) postMessage(message string, channel string, params slack.PostMessageParameters) {
	fmt.Printf("Posting message to slack:\n%v\n", message)
	params.Username = botName
	slackClient.SlackWrapper.PostMessage(channel, message, params)
}


func (slackClient *Slack) GetUserDetails(user string) (username, author string, ok bool) {
	if slackUser, err := slackClient.SlackWrapper.GetUserInfo(user); err == nil {
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