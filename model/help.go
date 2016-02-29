package model

type Help struct{ *Entry }

func NewHelp(clock Clock, author string) (help Help) {
	help = Help{NewEntry(clock, author)}
	return
}

func (help Help) String() string {
	return "helps" + help.Entry.String()
}

func (help Help) MakeCreateRequest() (request WhiteboardRequest) {
	request = help.Entry.MakeCreateRequest()
	request.Item.Kind = "Help"
	return
}

func (help Help) MakeUpdateRequest() (request WhiteboardRequest) {
	request = help.Entry.MakeUpdateRequest()
	request.Item.Kind = "Help"
	return
}
