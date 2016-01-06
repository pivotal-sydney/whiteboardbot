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
	StandupId int           `json:"-"`
}

type StandupItems struct {
	Helps        []Entry            `json:"Help"`
	Interestings []Entry  			`json:"Interesting"`
	Faces        []Entry            `json:"New face"`
	Events       []Entry        	`json:"Event"`
}

func NewEntry(clock Clock, author, title string, standup Standup) *Entry {
	location, err := time.LoadLocation(standup.TimeZone)
	if err != nil {
		location = time.Local
	}
	return &Entry{Date: clock.Now().In(location).Format(DATE_FORMAT), Author: author, Title: title, StandupId: standup.Id}
}

func (entry Entry) Validate() bool {
	return entry.Title != ""
}

func (entry Entry) MakeCreateRequest() WhiteboardRequest {
	return WhiteboardRequest{Token: os.Getenv("WB_AUTH_TOKEN"), Item: createItem(entry), Commit: "Create Item"}
}

func (entry Entry) MakeUpdateRequest() WhiteboardRequest {
	return WhiteboardRequest{Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: createItem(entry), Commit: "Update Item", Id: entry.Id}
}

func (entry Entry) GetEntry() *Entry {
	return &entry
}

func (entry Entry) String() string {
	if len(entry.Body) == 0 {
		return fmt.Sprintf("*%v*\n%v", entry.Title, entry.GetDateString())
	} else {
		return fmt.Sprintf("*%v*\n%v\n%v", entry.Title, entry.Body, entry.GetDateString())
	}
}

func (entry Entry) GetDateString() string {
	date, err := time.Parse(DATE_FORMAT, entry.Date)
	if err != nil {
		return entry.Date
	}
	return date.Format(DATE_STRING_FORMAT)
}

func createItem(entry Entry) Item {
	return Item{StandupId: entry.StandupId, Title: slackUnescape(entry.Title), Date: entry.Date, Public: "false", Description: slackUnescape(entry.Body), Author: entry.Author}
}

func (items StandupItems) FacesString() string {
	return toString("NEW FACES", items.Faces)
}

func (items StandupItems) InterestingsString() string {
	return toString("INTERESTINGS", items.Interestings)
}

func (items StandupItems) HelpsString() string {
	return toString("HELPS", items.Helps)
}

func (items StandupItems) EventsString() string {
	return toString("EVENTS", items.Events)
}

func toString(typeName string, entries []Entry) string {
	var buffer bytes.Buffer
	buffer.WriteString(typeName + "\n\n")
	for _, entry := range entries {
		buffer.WriteString(entry.String() + "\n \n")
	}
	return strings.TrimSuffix(buffer.String(), "\n \n")
}

func (items StandupItems) String() string {
	return fmt.Sprintf(">>>— — —\n \n \n \n%v\n \n \n \n%v\n \n \n \n%v\n \n \n \n%v\n \n \n \n— — —\n:clap:", items.FacesString(), items.HelpsString(), items.InterestingsString(), items.EventsString())
}

func (items StandupItems) Empty() bool {
	return len(items.Faces) == 0 && len(items.Events) == 0 && len(items.Helps) == 0 && len(items.Interestings) == 0
}

func slackUnescape(escaped string) string {
	return strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">").Replace(escaped)
}