package model

type Interesting struct{ *Entry }

func NewInteresting(clock Clock, author, title string, standup Standup) interface{} {
	return Interesting{NewEntry(clock, author, title, standup)}
}

func (interesting Interesting) MakeCreateRequest() (request WhiteboardRequest) {
	request = interesting.Entry.MakeCreateRequest()
	request.Item.Kind = "Interesting"
	return
}

func (interesting Interesting) MakeUpdateRequest() (request WhiteboardRequest) {
	request = interesting.Entry.MakeUpdateRequest()
	request.Item.Kind = "Interesting"
	return
}

func (interesting Interesting) GetEntry() *Entry {
	return interesting.Entry
}