package model
type WhiteboardRequest struct {
	Utf8   string `json:"utf8"`
	Method string `json:"_method,omitempty"`
	Token  string `json:"authenticity_token"`
	Item   Item `json:"item"`
	Commit string `json:"commit,omitempty"`
	Id     string `json:"id,omitempty"`
}

type Item struct {
	StandupId   int `json:"standup_id"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	PostId      string `json:"post_id"`
	Public      string `json:"public"`
	Kind        string `json:"kind"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
}
