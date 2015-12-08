package entry
import (
	"time"
	"fmt"
)

type EntryType interface {
	Validate() bool
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
