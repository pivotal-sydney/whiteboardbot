package model
import (
	"os"
	"fmt"
	"time"
	"strings"
)

const (
	DATE_FORMAT = "2006-01-02"
	DATE_STRING_FORMAT = "02 Jan 2006"
)

type EntryType interface {
	Validate() bool
	MakeCreateRequest() WhiteboardRequest
	MakeUpdateRequest() WhiteboardRequest
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
	ItemKind  string        `json:"-"`
}

func NewEntry(clock Clock, author, title string, standup Standup, itemKind string) *Entry {
	location, err := time.LoadLocation(standup.TimeZone)
	if err != nil {
		location = time.Local
	}
	return &Entry{Date: clock.Now().In(location).Format(DATE_FORMAT), Author: author, Title: title, StandupId: standup.Id, ItemKind: itemKind}
}

func (entry Entry) Validate() bool {
	return entry.Title != ""
}

func (entry Entry) MakeCreateRequest() WhiteboardRequest {
	return WhiteboardRequest{Token: os.Getenv("WB_AUTH_TOKEN"), Item: entry.toItem(), Commit: "Create Item"}
}

func (entry Entry) MakeUpdateRequest() WhiteboardRequest {
	return WhiteboardRequest{Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: entry.toItem(), Commit: "Update Item", Id: entry.Id}
}

func (entry Entry) GetEntry() *Entry {
	return &entry
}

func (entry Entry) String() string {
	author_str := ""
	if len(entry.Author) != 0 {
		author_str = fmt.Sprintf("\n[%v]", entry.Author)
	}
	body_str := ""
	if len(entry.Body) != 0 {
		body_str = fmt.Sprintf("\n%v", entry.Body)
	}
	return fmt.Sprintf("*%v*%v%v\n%v", entry.Title, body_str, author_str, entry.GetDateString())
}

func (entry Entry) GetDateString() string {
	date, err := time.Parse(DATE_FORMAT, entry.Date)
	if err != nil {
		return entry.Date
	}
	return date.Format(DATE_STRING_FORMAT)
}

func (entry Entry) toItem() Item {
	return Item{StandupId: entry.StandupId, Title: slackUnescape(entry.Title), Date: entry.Date, Public: "false", Description: slackUnescape(entry.Body), Author: entry.Author, Kind: entry.ItemKind}
}

func slackUnescape(escaped string) string {
	return strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">").Replace(escaped)
}
