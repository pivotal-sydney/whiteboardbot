package model

type Help struct{ *Entry }

func NewHelp(clock Clock, author, title string, standup Standup) EntryType {
	return Help{NewEntry(clock, author, title, standup, "Help")}
}

func (help Help) GetEntry() *Entry {
	return help.Entry
}
