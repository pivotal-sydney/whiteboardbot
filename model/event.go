package model
import (
	"fmt"
	"os"
)

type Event struct { *Entry }

func NewEvent(clock Clock, author string) (event Event) {
	event = Event{NewEntry(clock, author)}
	return
}

func (event Event) String() string {
	return fmt.Sprintf("events\n  *title: %v\n  body: %v\n  date: %v", event.Title, event.Body, event.Time.Format("2006-01-02"))
}

func (event Event) MakeCreateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: event.Title, Date: event.Time.Format("2006-01-02"), Public: "false", Kind: "Event", Description: event.Body, Author: event.Author}
	request = WhiteboardRequest{Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Create Item"}
	return
}

func  (event Event) MakeUpdateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: event.Title, Date: event.Time.Format("2006-01-02"), Public: "false", Kind: "Event", Description: event.Body, Author: event.Author}
	request = WhiteboardRequest{ Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Update Item", Id: event.Id}
	return
}