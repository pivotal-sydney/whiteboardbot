package model

type Interesting struct{ *Entry }

func NewInteresting(clock Clock, author, title string, standup Standup) (interesting interface{}) {
	interesting = Interesting{NewEntry(clock, author, title, standup)}
	return
}

func (interesting Interesting) String() string {
	return "INTERESTING" + interesting.Entry.String()
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