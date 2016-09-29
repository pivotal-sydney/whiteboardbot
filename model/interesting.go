package model

type Interesting struct{ *Entry }

func NewInteresting(clock Clock, author, title string, standup Standup) EntryType {
	return Interesting{NewEntry(clock, author, title, standup, "Interesting")}
}

func (interesting Interesting) GetEntry() *Entry {
	return interesting.Entry
}
