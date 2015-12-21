package model
import (
	"os"
	"fmt"
	"bytes"
	"time"
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

func NewEntry(clock Clock, author, title string, standup Standup) (entry *Entry) {
	location, err := time.LoadLocation(standup.TimeZone)
	if err != nil {
		location = time.Local
	}
	entry = &Entry{Date: clock.Now().In(location).Format("2006-01-02"), Author: author, Title: title, StandupId: standup.Id}
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

func (items StandupItems) FacesString() string {
	var buffer bytes.Buffer
	buffer.WriteString("New Faces:")
	for _, face := range items.Faces {
		buffer.WriteString("\n" + Face{&face}.String())
	}
	return buffer.String()
}

func (items StandupItems) InterestingsString() string {
	var buffer bytes.Buffer
	buffer.WriteString("Interestings:")
	for _, interesting := range items.Interestings {
		buffer.WriteString("\n" + Interesting{&interesting}.String())
	}
	return buffer.String()
}


func (items StandupItems) HelpsString() string {
	var buffer bytes.Buffer
	buffer.WriteString("Helps:")
	for _, help := range items.Helps {
		buffer.WriteString("\n" + Help{&help}.String())
	}
	return buffer.String()
}

func (items StandupItems) EventsString() string {
	var buffer bytes.Buffer
	buffer.WriteString("Events:")
	for _, event := range items.Events {
		buffer.WriteString("\n" + Event{&event}.String())
	}
	return buffer.String()
}

func (items StandupItems) String() string {
	return fmt.Sprintf("%v\n%v\n%v\n%v", items.FacesString(), items.InterestingsString(), items.HelpsString(), items.EventsString())
}

func (items StandupItems) Empty() bool {
	return len(items.Faces) == 0 && len(items.Events) == 0 && len(items.Helps) == 0 && len(items.Interestings) == 0
}
