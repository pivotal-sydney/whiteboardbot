package app

import (
	"errors"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"strconv"
)

type StandupRepository interface {
	FindStandup(string) (Standup, error)
	SaveEntry(EntryType) (PostResult, error)
	GetStandupItems(string) (StandupItems, error)
}

type PostResult struct {
	ItemId string
}

type WhiteboardGateway struct {
	RestClient RestClient
}

func (gateway WhiteboardGateway) FindStandup(standupId string) (standup Standup, err error) {
	standup, ok := gateway.RestClient.GetStandup(standupId)
	if !ok {
		err = errors.New("Standup not found!")
	}
	return
}

func (gateway WhiteboardGateway) SaveEntry(entryType EntryType) (PostResult, error) {
	itemId, ok := PostEntryToWhiteboard(gateway.RestClient, entryType)

	if !ok {
		return PostResult{}, errors.New("Problem creating post.")
	}

	return PostResult{itemId}, nil
}

func (gateway WhiteboardGateway) GetStandupItems(standupId string) (standupItems StandupItems, err error) {
	standupIdInt, err := strconv.Atoi(standupId)
	if err != nil {
		err = errors.New(MISSING_STANDUP)
		return
	}

	standupItems, ok := gateway.RestClient.GetStandupItems(standupIdInt)
	if !ok {
		err = errors.New("Failed fetching standup items")
	}

	return
}
