package model
import (
	"fmt"
	"os"
)

type Interesting struct { *Entry }

func NewInteresting(clock Clock, author string) (interesting Interesting) {
	interesting = Interesting{NewEntry(clock, author)}
	return
}

func (interesting Interesting) String() string {
	return fmt.Sprintf("interestings\n  *title: %v\n  body: %v\n  date: %v", interesting.Title, interesting.Body, interesting.Time.Format("2006-01-02"))
}

func (interesting Interesting) MakeCreateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: interesting.Title, Date: interesting.Time.Format("2006-01-02"), Public: "false", Kind: "Interesting", Description: interesting.Body, Author: interesting.Author}
	request = WhiteboardRequest{Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Create Item"}
	return
}

func  (interesting Interesting) MakeUpdateRequest() (request WhiteboardRequest) {
	item := Item{StandupId: 1, Title: interesting.Title, Date: interesting.Time.Format("2006-01-02"), Public: "false", Kind: "Interesting", Description: interesting.Body, Author: interesting.Author}
	request = WhiteboardRequest{ Method: "patch", Token: os.Getenv("WB_AUTH_TOKEN"), Item: item, Commit: "Update Item", Id: interesting.Id}
	return
}


