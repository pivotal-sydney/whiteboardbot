package model
import (
	"time"
	"os"
	"fmt"
)

type EntryType interface {
	Validate() bool
	MakeCreateRequest() (request WhiteboardRequest)
	MakeUpdateRequest() (request WhiteboardRequest)
	String() string
}

type Entry struct {
	Time   time.Time
	Title  string
	Body   string
	Author string
	Id     string
}

func NewEntry(clock Clock, author string) (entry *Entry) {
	entry = &Entry{Time: clock.Now(), Author: author}
	return
}

func (entry Entry) Validate() bool {
	return entry.Title != ""
}

func (entry Entry) MakeCreateRequest() (request WhiteboardRequest) {
	request = WhiteboardRequest{Token: os.Getenv("WB_AUTH_TOKEN"), Item: createItem(entry), Commit: "Create Item"}
	return
}

func (entry Entry) MakeUpdateRequest() (request WhiteboardRequest) {
	request = WhiteboardRequest{ Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: createItem(entry), Commit: "Update Item", Id: entry.Id}
	return
}

func (entry Entry) String() string {
	return fmt.Sprintf("\n  *title: %v\n  body: %v\n  date: %v", entry.Title, entry.Body, entry.Time.Format("2006-01-02"))
}

func createItem(entry Entry) (item Item) {
	item = Item{StandupId: 1, Title: entry.Title, Date: entry.Time.Format("2006-01-02"), Public: "false", Description: entry.Body, Author: entry.Author}
	return
}