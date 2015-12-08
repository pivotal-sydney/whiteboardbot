package model
import (
	"time"
	"fmt"
	"os"
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

type Face struct { Entry }
type Interesting struct { Entry }

func NewEntry(clock Clock, author string) (entry Entry) {
	entry = Entry{Time: clock.Now(), Author: author}
	return
}

func NewFace(clock Clock, author string) (face Face) {
	face = Face{NewEntry(clock, author)}
	return
}

func NewInteresting(clock Clock, author string) (interesting Interesting) {
	interesting = Interesting{NewEntry(clock, author)}
	return
}

func (entry Entry) Validate() bool {
	return entry.Title != ""
}

func (entry Entry) String() string {
	return fmt.Sprintf("faces\n  *name: %v\n  date: %v", entry.Title, entry.Time.Format("2006-01-02"))
}

func (interesting Interesting) String() string {
	return fmt.Sprintf("interestings\n  *title: %v\n  body: %v\n  date: %v", interesting.Title, interesting.Body, interesting.Time.Format("2006-01-02"))
}

func (face Face) MakeCreateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: face.Title, Date: face.Time.Format("2006-01-02"), Public: "false", Kind: "New face", Author: face.Author}
	request = WhiteboardRequest{Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Create New Face"}
	return
}

func  (face Face) MakeUpdateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: face.Title, Date: face.Time.Format("2006-01-02"), Public: "false", Kind: "New face", Author: face.Author}
	request = WhiteboardRequest{ Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Update New Face", Id: face.Id}
	return
}


