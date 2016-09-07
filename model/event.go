package model

type Event struct{ *Entry }

func NewEvent(clock Clock, author, title string, standup Standup) interface{} {
	return Event{NewEntry(clock, author, title, standup, "Event")}
}

func (event Event) GetEntry() *Entry {
	return event.Entry
}
