package app

import (
	"fmt"
	"github.com/nlopes/slack"
)

type Slack struct {
	SlackRtm *slack.RTM
}

type SlackUser struct {
	Username string
	Author   string
	TimeZone string
}

type SlackChannel struct {
	Id   string
	Name string
}

type SlackContext struct {
	User    SlackUser
	Channel SlackChannel
}

type SlackClient interface {
	PostMessage(message string, channel string)
	PostMessageWithMarkdown(message string, channel string)
	GetUserDetails(user string) (slackUser SlackUser)
	GetChannelDetails(channel string) (slackChannel SlackChannel)
}

func (slackClient *Slack) PostMessage(message string, channel string) {
	slackClient.postMessage(message, channel, slack.PostMessageParameters{})
}

func (slackClient *Slack) PostMessageWithMarkdown(message string, channel string) {
	slackClient.postMessage(message, channel, slack.PostMessageParameters{Markdown: true})
}

func (slackClient *Slack) postMessage(message string, channel string, params slack.PostMessageParameters) {
	fmt.Printf("Posting message to slack:\n%v\n", message)
	params.AsUser = true
	slackClient.SlackRtm.PostMessage(channel, message, params)
}

func (slackClient *Slack) GetUserDetails(user string) (slackUser SlackUser) {
	if userInfo, err := slackClient.SlackRtm.GetUserInfo(user); err == nil {
		slackUser.Username = userInfo.Name
		slackUser.Author = GetAuthor(userInfo)
		slackUser.TimeZone = userInfo.TZ
	} else {
		slackUser.Username = user
		slackUser.Author = user
		slackUser.TimeZone = "America/Los_Angeles"
		fmt.Printf("SlackClient.GetUserDetails returned error: %v, %v\n", user, err)
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

func (slackClient *Slack) GetChannelDetails(channel string) SlackChannel {
	slackChannel, err := slackClient.SlackRtm.GetChannelInfo(channel)
	if err != nil {
		slackChannel = &slack.Channel{}
		slackChannel.ID = channel
		slackChannel.Name = "unknown"
	}
	return SlackChannel{Id: slackChannel.ID, Name: slackChannel.Name}
}

func handleMissingEntry(slackClient SlackClient, channel string) {
	slackClient.PostMessageWithMarkdown(THUMBS_DOWN+"Hey, you forgot to start new entry. Start with one of `wb [face interesting help event] [title]` first!", channel)
}

func handleNotRegistered(slackClient SlackClient, channel string) {
	slackClient.PostMessage(THUMBS_DOWN+"You haven't registered your standup yet. wb r <id> first!", channel)
	return
}

func handleStandupNotFound(slackClient SlackClient, standupId string, channel string) {
	slackClient.PostMessage(fmt.Sprintf("%sI couldn't find a standup with id: %v", THUMBS_DOWN, standupId), channel)
	return
}
