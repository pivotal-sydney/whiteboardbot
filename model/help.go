package model
import (
	"fmt"
	"os"
)

type Help struct { *Entry }

func NewHelp(clock Clock, author string) (help Help) {
	help = Help{NewEntry(clock, author)}
	return
}

func (help Help) String() string {
	return fmt.Sprintf("helps\n  *title: %v\n  body: %v\n  date: %v", help.Title, help.Body, help.Time.Format("2006-01-02"))
}

func (help Help) MakeCreateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: help.Title, Date: help.Time.Format("2006-01-02"), Public: "false", Kind: "Help", Description: help.Body, Author: help.Author}
	request = WhiteboardRequest{Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Create Item"}
	return
}

func  (help Help) MakeUpdateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: help.Title, Date: help.Time.Format("2006-01-02"), Public: "false", Kind: "Help", Description: help.Body, Author: help.Author}
	request = WhiteboardRequest{ Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Update Item", Id: help.Id}
	return
}