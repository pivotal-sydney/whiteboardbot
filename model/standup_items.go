package model
import (
	"fmt"
	"strings"
	"bytes"
)

type StandupItems struct {
	Helps        []Entry            `json:"Help"`
	Interestings []Entry  			`json:"Interesting"`
	Faces        []Entry            `json:"New face"`
	Events       []Entry        	`json:"Event"`
}

func (items StandupItems) FacesString() string {
	return toString("NEW FACES", items.Faces)
}

func (items StandupItems) InterestingsString() string {
	return toString("INTERESTINGS", items.Interestings)
}

func (items StandupItems) HelpsString() string {
	return toString("HELPS", items.Helps)
}

func (items StandupItems) EventsString() string {
	return toString("EVENTS", items.Events)
}

func toString(typeName string, entries []Entry) string {
	var buffer bytes.Buffer
	buffer.WriteString(typeName + "\n\n")
	for _, entry := range entries {
		buffer.WriteString(entry.String() + "\n \n")
	}
	return strings.TrimSuffix(buffer.String(), "\n \n")
}

func (items StandupItems) String() string {
	return fmt.Sprintf(">>>— — —\n \n \n \n%v\n \n \n \n%v\n \n \n \n%v\n \n \n \n%v\n \n \n \n— — —\n:clap:", items.FacesString(), items.HelpsString(), items.InterestingsString(), items.EventsString())
}

func (items StandupItems) Empty() bool {
	return len(items.Faces) == 0 && len(items.Events) == 0 && len(items.Helps) == 0 && len(items.Interestings) == 0
}