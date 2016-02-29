package model

type Event struct{ *Entry }

func NewEvent(clock Clock, author string) (event Event) {
	event = Event{NewEntry(clock, author)}
	return
}

func (event Event) String() string {
	return "events" + event.Entry.String()
}

func (event Event) MakeCreateRequest() (request WhiteboardRequest) {
	request = event.Entry.MakeCreateRequest()
	request.Item.Kind = "Event"
	return
}

func (event Event) MakeUpdateRequest() (request WhiteboardRequest) {
	request = event.Entry.MakeUpdateRequest()
	request.Item.Kind = "Event"
	return
}
