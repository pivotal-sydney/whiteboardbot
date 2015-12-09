package model
import (
	"os"
	"fmt"
)

type Face struct { *Entry }

func NewFace(clock Clock, author string) (face Face) {
	face = Face{NewEntry(clock, author)}
	return
}

func (face Face) String() string {
	return fmt.Sprintf("faces\n  *name: %v\n  date: %v", face.Title, face.Time.Format("2006-01-02"))
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

