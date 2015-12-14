package model

type Interesting struct{ *Entry }

func NewInteresting(clock Clock, author, title string) (interesting Interesting) {
	interesting = Interesting{NewEntry(clock, author, title)}
	return
}

func (interesting Interesting) String() string {
	return "interestings" + interesting.Entry.String()
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