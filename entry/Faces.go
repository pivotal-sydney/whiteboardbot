package entry
import (
	"time"
	"fmt"
)

type Face struct {
	Time time.Time
	Name string
}

func NewFace(clock Clock) (face Face) {
	face = Face{}
	face.Time = clock.Now()
	return
}

func (face Face) GetDate() time.Time {
	return face.Time
}

func (face *Face) SetDate(time time.Time) {
	face.Time = time
}

func (face Face) GetName() string {
	return face.Name
}

func (face *Face) SetName(name string) {
	face.Name = name
}

func (face Face) Validate() bool {
	return face.Name != ""
}

func (face Face) String() string {
	return fmt.Sprintf("faces\n  *name: %v\n  date: %v", face.Name, face.Time.Format("2006-01-02"))
}
