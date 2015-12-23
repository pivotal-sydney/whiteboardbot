package spec

import (
	"github.com/xtreme-andleung/whiteboardbot/model"
	"time"
	"strconv"
)

type MockSlackClient struct {
	PostMessageCalled bool
	Message           string
	EntryType 		  model.EntryType
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

func (slackClient *MockSlackClient) PostEntry(entryType model.EntryType, channel string, status string) {
	slackClient.EntryType = entryType
	slackClient.Status = status
}

func (slackClient *MockSlackClient) GetUserDetails(user string) (username, author string, ok bool) {
	username = user
	if username == "" {
		username = "aleung"
	}
	author = "Andrew Leung"
	ok = true
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

func (client *MockRestClient) Post(request model.WhiteboardRequest, standupId int) (itemId string, ok bool) {
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
