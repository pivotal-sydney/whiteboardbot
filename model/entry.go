package model
import (
	"time"
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
	return
}

func (entry Entry) MakeUpdateRequest() (request WhiteboardRequest) {
	return
}
func (entry Entry) String() string {
	return ""
}