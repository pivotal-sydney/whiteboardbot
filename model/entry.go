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
}

type Entry struct {
	Time   time.Time
	Title  string
	Body   string
	Author string
	Id     string
}

type Face struct { Entry }

func NewFace(clock Clock, author string) (face Face) {
	face = Face{Entry{clock.Now(), "", "", author, ""}}
	return
}

func (entry Entry) Validate() bool {
	return entry.Title != ""
}

func (face Face) String() string {
	return fmt.Sprintf("faces\n  *name: %v\n  date: %v", face.Title, face.Time.Format("2006-01-02"))
}

func (entry Entry) String() string {
	return fmt.Sprintf("\n  *title: %v\n  body: %v\n  date: %v", entry.Title, entry.Body, entry.Time.Format("2006-01-02"))
}

func (face Face) MakeCreateRequest() (request WhiteboardRequest) {
	item := Item{1, face.Title, face.Time.Format("2006-01-02"), "", "false", "New face", "", face.Author}
	request = WhiteboardRequest{"", "", os.Getenv("WB_AUTH_TOKEN"), item, "Create New Face", ""}
	return
}

func  (face Face) MakeUpdateRequest() (request WhiteboardRequest) {
	item := Item{1, face.Title, face.Time.Format("2006-01-02"), "", "false", "New face", "", face.Author}
	request = WhiteboardRequest{"", "patch", os.Getenv("WB_AUTH_TOKEN"), item, "Update New Face", face.Id}
	return
}


