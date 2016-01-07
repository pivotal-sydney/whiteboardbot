package spec

import (
	"github.com/xtreme-andleung/whiteboardbot/model"
	"time"
	"strconv"
	"encoding/json"
	. "github.com/xtreme-andleung/whiteboardbot/app"
)

type MockSlackClient struct {
	PostMessageCalled bool
	Message           string
	Entry 		      *model.Entry
	Status 			  string
}

func (slackClient *MockSlackClient) PostMessage(message string, channel string, status string) {
	slackClient.PostMessageCalled = true
	slackClient.Message = message
	slackClient.Status = status
}

func (slackClient *MockSlackClient) PostMessageWithMarkdown(message string, channel string, status string) {
	slackClient.PostMessageCalled = true
	slackClient.Message = message
	slackClient.Status = status
}

func (slackClient *MockSlackClient) PostEntry(entry *model.Entry, channel string, status string) {
	slackClient.Entry = entry
	slackClient.Status = status
}

func (slackClient *MockSlackClient) GetUserDetails(user string) (slackUser SlackUser) {
	slackUser.Username = user
	if slackUser.Username == "" {
		slackUser.Username = "aleung"
	}
	slackUser.Author = "Andrew Leung"
	slackUser.TimeZone = "Australia/Sydney"
	return
}

type MockClock struct{}

func (clock MockClock) Now() time.Time {
	return time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)
}

type MockRestClient struct {
	PostCalledCount int
	Request         model.WhiteboardRequest
	StandupItems    model.StandupItems
}

func (client MockRestClient) GetStandupItems(standupId int) (items model.StandupItems, ok bool) {
	items = client.StandupItems
	ok = true
	return
}

func (client *MockRestClient) Post(request model.WhiteboardRequest) (itemId string, ok bool) {
	client.PostCalledCount++
	client.Request = request
	ok = true
	itemId = "1"
	return
}

func (*MockRestClient) GetStandup(standupId string) (standup model.Standup, ok bool) {
	id, _ := strconv.Atoi(standupId)
	standup.Id = id
	standup.TimeZone = "Australia/Sydney"
	standup.Title = "Sydney"
	ok = true
	return
}

type MockStore struct {
	StoreMap map[string]string
}

func (store *MockStore) Get(key string) (value string, ok bool) {
	if store.StoreMap == nil {
		store.StoreMap = make(map[string]string)
	}
	value, ok = store.StoreMap[key]
	return
}

func (store *MockStore) Set(key string, value string) {
	if store.StoreMap == nil {
		store.StoreMap = make(map[string]string)
	}
	store.StoreMap[key] = value
}

func (store *MockStore) GetStandup(channel string) (standup model.Standup, ok bool) {
	var standupJson string
	standupJson, _ = store.Get(channel)
	err := json.Unmarshal([]byte(standupJson), &standup)
	ok = err == nil
	return
}

func (store *MockStore) SetStandup(channel string, standup model.Standup) {
	standupJson, _ := json.Marshal(standup)
	store.Set(channel, string(standupJson))
}