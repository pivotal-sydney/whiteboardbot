package model
import (
	"os"
	"fmt"
)

type EntryType interface {
	Validate() bool
	MakeCreateRequest() (request WhiteboardRequest)
	MakeUpdateRequest() (request WhiteboardRequest)
	String() string
	GetEntry() *Entry
}

type Entry struct {
	Date      string     	`json:"date"`
	Title     string        `json:"title"`
	Body      string        `json:"description"`
	Author    string        `json:"author"`
	Id        string        `json:"-"`
	StandupId int            `json:"-"`
}

type StandupItems struct {
	Helps        []Entry            `json:"Help"`
	Interestings []Entry  			`json:"Interesting"`
	Faces        []Entry            `json:"New face"`
	Events       []Entry        	`json:"Event"`
}

func NewEntry(clock Clock, author, title string, standupId int) (entry *Entry) {
	entry = &Entry{Date: clock.Now().Format("2006-01-02"), Author: author, Title: title, StandupId: standupId}
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
	request = WhiteboardRequest{Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: createItem(entry), Commit: "Update Item", Id: entry.Id}
	return
}

func (entry Entry) GetEntry() *Entry {
	return &entry
}

func (entry Entry) String() string {
	return fmt.Sprintf("\n  *title: %v\n  body: %v\n  date: %v", entry.Title, entry.Body, entry.Date)
}

func createItem(entry Entry) (item Item) {
	item = Item{StandupId: entry.StandupId, Title: entry.Title, Date: entry.Date, Public: "false", Description: entry.Body, Author: entry.Author}
	return
}