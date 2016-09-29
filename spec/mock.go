package spec

import (
	"encoding/json"
	"errors"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	"github.com/pivotal-sydney/whiteboardbot/model"
	"strconv"
	"time"
)

type MockSlackClient struct {
	PostMessageCalled bool
	Message           string
	Entry             *model.Entry
	Status            string
	ChannelId         string
	SlackUserMap      map[string]SlackUser
	SlackChannelMap   map[string]SlackChannel
}

func (slackClient *MockSlackClient) PostMessage(message string, channel string) {
	slackClient.PostMessageCalled = true
	slackClient.Message = message
	slackClient.ChannelId = channel
}

func (slackClient *MockSlackClient) PostMessageWithMarkdown(message string, channel string) {
	slackClient.PostMessageCalled = true
	slackClient.Message = message
	slackClient.ChannelId = channel
}

func (slackClient *MockSlackClient) PostEntry(entry *model.Entry, channel string, status string) {
	slackClient.Entry = entry
	slackClient.Status = status
	slackClient.ChannelId = channel
}

func (slackClient *MockSlackClient) GetUserDetails(user string) (slackUser SlackUser) {
	slackClient.initSlackUserMap()
	slackUser, ok := slackClient.SlackUserMap[user]
	if !ok {
		slackUser.Username = user
		slackUser.Author = user
		slackUser.TimeZone = "America/Los_Angeles"
	}
	return
}

func (slackClient *MockSlackClient) initSlackUserMap() {
	if slackClient.SlackUserMap == nil {
		slackClient.SlackUserMap = map[string]SlackUser{
			"U987": {
				Username: "aleung",
				Author:   "Andrew Leung",
				TimeZone: "Australia/Sydney",
			},
			"UUserId": {
				Username: "user-name",
				Author:   "Andrew Leung",
				TimeZone: "Australia/Sydney",
			},
			"UUserId2": {
				Username: "user-name-two",
				Author:   "Andrew Leung",
				TimeZone: "Australia/Sydney",
			},
		}
	}
}

func (slackClient *MockSlackClient) AddSlackUser(userId string, user SlackUser) {
	slackClient.initSlackUserMap()
	slackClient.SlackUserMap[userId] = user
}

func (slackClient *MockSlackClient) GetChannelDetails(channelId string) (slackChannel SlackChannel) {
	slackClient.initSlackChannelMap()
	slackChannel, ok := slackClient.SlackChannelMap[channelId]
	if !ok {
		slackChannel.Id = channelId
		slackChannel.Name = "unknown"
	}

	return
}

func (slackClient *MockSlackClient) initSlackChannelMap() {
	if slackClient.SlackChannelMap == nil {
		channel1 := SlackChannel{Id: "CChannelId", Name: "channel-name"}
		channel2 := SlackChannel{Id: "CChannelId2", Name: "channel-name-two"}

		slackClient.SlackChannelMap = map[string]SlackChannel{
			channel1.Id: channel1,
			channel2.Id: channel2,
		}
	}
}

func (slackClient *MockSlackClient) AddSlackChannel(channelId string, channel SlackChannel) {
	slackClient.initSlackChannelMap()
	slackClient.SlackChannelMap[channelId] = channel
}

type MockClock struct{}

func (clock MockClock) Now() time.Time {
	return time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)
}

type MockRestClient struct {
	PostCalledCount     int
	Request             model.WhiteboardRequest
	StandupItems        model.StandupItems
	StandupMap          map[int]model.Standup
	failPost            bool
	postItemId          string
	failGetStandupItems bool
}

func (client MockRestClient) GetStandupItems(standupId int) (items model.StandupItems, ok bool) {
	if client.failGetStandupItems {
		ok = false
	} else {
		items = client.StandupItems
		ok = true
	}
	return
}

func (client *MockRestClient) SetGetStandupItemsError() {
	client.failGetStandupItems = true
}

func (client *MockRestClient) Post(request model.WhiteboardRequest) (itemId string, ok bool) {
	client.PostCalledCount++
	if !client.failPost {
		client.Request = request
		ok = true
		itemId = client.getNextId()
	}
	return
}

func (client *MockRestClient) getNextId() (id string) {
	id = client.postItemId
	if id == "" {
		id = "1"
	}

	client.postItemId += client.postItemId
	return
}

func (client *MockRestClient) SetPostError() {
	client.failPost = true
}

func (client *MockRestClient) SetPostItemId(id string) {
	client.postItemId = id
}

func (client *MockRestClient) SetStandup(standup model.Standup) {
	if client.StandupMap == nil {
		client.StandupMap = make(map[int]model.Standup)
	}
	client.StandupMap[standup.Id] = standup
}

func (client *MockRestClient) GetStandup(standupId string) (standup model.Standup, ok bool) {
	id, _ := strconv.Atoi(standupId)
	standup, ok = client.StandupMap[id]
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

type MockQuietWhiteboard struct {
	HandleInputCalled bool
	HandleInputArgs   struct {
		Text    string
		Context SlackContext
	}
}

func (mqw *MockQuietWhiteboard) ProcessCommand(input string, context SlackContext) CommandResult {
	mqw.HandleInputCalled = true
	mqw.HandleInputArgs.Text = input
	mqw.HandleInputArgs.Context = context

	return &MessageCommandResult{Text: "This is a mock message"}
}

type MockWhiteboardGateway struct {
	StandupMap          map[int]model.Standup
	SaveEntryCalled     bool
	EntrySaved          model.EntryType
	GetStandupId        string
	failSaveEntry       bool
	failGetStandupItems bool
}

func (gateway *MockWhiteboardGateway) FindStandup(standupId string) (standup model.Standup, err error) {
	var ok bool
	id, _ := strconv.Atoi(standupId)

	standup, ok = gateway.StandupMap[id]

	if !ok {
		err = errors.New("Standup not found!")
	}

	return
}

func (gateway *MockWhiteboardGateway) SaveEntry(entryType model.EntryType) (PostResult, error) {
	if gateway.failSaveEntry {
		return PostResult{}, errors.New("Problem creating post.")
	}
	gateway.SaveEntryCalled = true
	gateway.EntrySaved = entryType

	return PostResult{ItemId: "1"}, nil
}

func (gateway *MockWhiteboardGateway) SetSaveEntryError() {
	gateway.failSaveEntry = true
}

func (gateway *MockWhiteboardGateway) SetStandup(standup model.Standup) {
	if gateway.StandupMap == nil {
		gateway.StandupMap = make(map[int]model.Standup)
	}
	gateway.StandupMap[standup.Id] = standup
}

func (gateway *MockWhiteboardGateway) GetStandupItems(standupId string) (standupItems model.StandupItems, err error) {
	if gateway.failGetStandupItems {
		err = errors.New("Error retrieving standup items.")
	} else {
		standupItems = model.StandupItems{Interestings: []model.Entry{
			{Title: "Interesting 1", Author: "Alice", Date: "2015-01-02"},
			{Title: "Interesting 2", Author: "Bob", Date: "2015-01-12"},
		}}
	}
	return standupItems, err
}

func (gateway *MockWhiteboardGateway) SetGetStandupItemsError() {
	gateway.failGetStandupItems = true
}
