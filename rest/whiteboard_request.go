package rest
import (
	"os"
	. "github.com/xtreme-andleung/whiteboardbot/entry"
)
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
	Public      bool `json:"public"`
	Kind        string `json:"kind"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
}

type FaceRequest WhiteboardRequest

func NewCreateFaceRequest(face *Face) (*FaceRequest) {
	item := Item{1, face.Name, face.Time.Format("2006-01-02"), "", false, "New Face", "", face.Author}
	return &FaceRequest{"", "", os.Getenv("token"), item, "Create New Face", ""}
}