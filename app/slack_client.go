package app
import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/model"
)

type Slack struct {
	SlackWrapper SlackWrapper
}

type SlackClient interface {
	PostMessage(message string, channel string, status string)
	PostMessageWithMarkdown(message string, channel string, status string)
	PostEntry(entryType model.EntryType, channel string, status string)
	GetUserDetails(user string) (username, author string)
}

func (slackClient *Slack) PostMessage(message string, channel string, status string) {
	slackClient.postMessage(message, channel, status, slack.PostMessageParameters{})
}

func (slackClient *Slack) PostMessageWithMarkdown(message string, channel string, status string) {
	slackClient.postMessage(message, channel, status, slack.PostMessageParameters{Markdown: true})
}

func (slackClient *Slack) PostEntry(entryType model.EntryType, channel string, status string) {
	message := entryType.String()
	slackClient.PostMessage(message, channel, status)
}

func (slackClient *Slack) postMessage(message string, channel string, status string, params slack.PostMessageParameters) {
	message = status + message
	fmt.Printf("Posting message to slack:\n%v\n", message)
	params.Username = BOT_NAME
	slackClient.SlackWrapper.PostMessage(channel, message, params)
}


func (slackClient *Slack) GetUserDetails(user string) (username, author string) {
	if slackUser, err := slackClient.SlackWrapper.GetUserInfo(user); err == nil {
		username = slackUser.Name
		author = GetAuthor(slackUser)
	} else {
		username = user
		author = user
		fmt.Printf("SlackClient.GetUserDetails returned error: %v, %v\n", username, err)
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

func handleMissingEntry(slackClient SlackClient, channel string) {
	slackClient.PostMessageWithMarkdown("Hey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!", channel, THUMBS_DOWN)
}

func handleNotRegistered(slackClient SlackClient, channel string) {
	slackClient.PostMessage("You haven't registered your standup yet. wb r <id> first!", channel, THUMBS_DOWN)
	return
}

func handleStandupNotFound(slackClient SlackClient, standupId string, channel string) {
	slackClient.PostMessage(fmt.Sprintf("I couldn't find a standup with id: %v", standupId), channel, THUMBS_DOWN)
	return
}