package model
import (
	"os"
	"fmt"
	"bytes"
	"time"
	"strings"
)

const (
	DATE_FORMAT = "2006-01-02"
	DATE_STRING_FORMAT = "02 Jan 2006"
)

type EntryType interface {
	Validate() bool
	MakeCreateRequest() (request WhiteboardRequest)
	MakeUpdateRequest() (request WhiteboardRequest)
	String() string
	GetEntry() *Entry
	GetDateString() string
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
	entry = &Entry{Date: clock.Now().In(location).Format(DATE_FORMAT), Author: author, Title: title, StandupId: standup.Id}
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
	return fmt.Sprintf("\n\n>*%v*\n>%v\n>%v", entry.Title, entry.Body, entry.GetDateString())
}

func (entry Entry) GetDateString() string {
	date, err := time.Parse(DATE_FORMAT, entry.Date)
	if err != nil {
		return entry.Date
	}
	return date.Format(DATE_STRING_FORMAT)
}

func createItem(entry Entry) (item Item) {
	item = Item{StandupId: entry.StandupId, Title: entry.Title, Date: entry.Date, Public: "false", Description: entry.Body, Author: entry.Author}
	return
}

func (items StandupItems) FacesString() string {
	var buffer bytes.Buffer
	for _, face := range items.Faces {
		buffer.WriteString(Face{&face}.String() + "\n")
	}
	return strings.TrimSuffix(buffer.String(), "\n")
}

func (items StandupItems) InterestingsString() string {
	var buffer bytes.Buffer
	for _, interesting := range items.Interestings {
		buffer.WriteString(Interesting{&interesting}.String() + "\n")
	}
	return strings.TrimSuffix(buffer.String(), "\n")
}


func (items StandupItems) HelpsString() string {
	var buffer bytes.Buffer
	for _, help := range items.Helps {
		buffer.WriteString(Help{&help}.String() + "\n")
	}
	return strings.TrimSuffix(buffer.String(), "\n")
}

func (items StandupItems) EventsString() string {
	var buffer bytes.Buffer
	for _, event := range items.Events {
		buffer.WriteString(Event{&event}.String() + "\n")
	}
	return strings.TrimSuffix(buffer.String(), "\n")
}

func (items StandupItems) String() string {
	return fmt.Sprintf("%v\n%v\n%v\n%v", items.FacesString(), items.InterestingsString(), items.HelpsString(), items.EventsString())
}

func (items StandupItems) Empty() bool {
	return len(items.Faces) == 0 && len(items.Events) == 0 && len(items.Helps) == 0 && len(items.Interestings) == 0
}
