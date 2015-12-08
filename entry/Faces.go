package entry
import (
	"time"
	"fmt"
)

type Face struct {
	Time time.Time
	Name string
	Author string
	Id string
}

func NewFace(clock Clock, author string) (face Face) {
	face = Face{clock.Now(), "", author, ""}
	return
}

func (face Face) Validate() bool {
	return face.Name != ""
}

func (face Face) String() string {
	return fmt.Sprintf("faces\n  *name: %v\n  date: %v", face.Name, face.Time.Format("2006-01-02"))
}
