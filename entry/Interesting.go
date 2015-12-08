package entry
import (
	"time"
	"fmt"
)

type Interesting struct {
	Time time.Time
	Title string
	Body string
	Author string
	Id string
}

func NewInteresting(clock Clock, author string) (interesting Interesting) {
	interesting = Interesting{clock.Now(), "", "", author, ""}
	return
}

func (interesting Interesting) Validate() bool {
	return interesting.Title != ""
}

func (interesting Interesting) String() string {
	return fmt.Sprintf("interestings\n  *title: %v\n  body: %v\n  date: %v", interesting.Title, interesting.Body, interesting.Time.Format("2006-01-02"))
}