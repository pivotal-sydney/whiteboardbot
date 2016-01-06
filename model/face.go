package model

type Face struct{ *Entry }

func NewFace(clock Clock, author, title string, standup Standup) interface{} {
	return Face{NewEntry(clock, author, title, standup, "New face")}
}

func (face Face) MakeCreateRequest() (request WhiteboardRequest) {
	request = face.Entry.MakeCreateRequest()
	request.Commit = "Create New Face"
	return
}

func (face Face) MakeUpdateRequest() (request WhiteboardRequest) {
	request = face.Entry.MakeUpdateRequest()
	request.Commit = "Update New Face"
	return
}

func (face Face) GetEntry() *Entry {
	return face.Entry
}