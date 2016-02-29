package model

type Interesting struct{ *Entry }

func NewInteresting(clock Clock, author string) (interesting Interesting) {
	interesting = Interesting{NewEntry(clock, author)}
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
