package entry
import (
	"time"
	"fmt"
)

type Face struct {
	Time time.Time
	Name string
	Author string
}

func NewFace(clock Clock, author string) (*Face) {
	face := Face{}
	face.Time = clock.Now()
	face.Author = author
	return &face
}

func (face Face) Validate() bool {
	return face.Name != ""
}

func (face Face) String() string {
	return fmt.Sprintf("faces\n  *name: %v\n  date: %v", face.Name, face.Time.Format("2006-01-02"))
}
