package model

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type StandupItems struct {
	Helps        []Entry `json:"Help"`
	Interestings []Entry `json:"Interesting"`
	Faces        []Entry `json:"New face"`
	Events       []Entry `json:"Event"`
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

func (items StandupItems) Filter(numberOfDays int, clock Clock, userTimeZone string) (newItems StandupItems) {
	location, err := time.LoadLocation(userTimeZone)
	if err != nil {
		location = time.Local
	}

	cutOff := clock.Now().In(location).AddDate(0, 0, numberOfDays)

	filterEntries := func(entries []Entry, cutOff time.Time, location *time.Location) (newEntries []Entry) {
		for _, entry := range entries {
			entryDate, _ := time.Parse("2006-01-02", entry.Date)
			if entryDate.In(location).Before(cutOff) {
				newEntries = append(newEntries, entry)
			}
		}

		return
	}

	newItems.Faces = filterEntries(items.Faces, cutOff, location)
	newItems.Helps = filterEntries(items.Helps, cutOff, location)
	newItems.Interestings = filterEntries(items.Interestings, cutOff, location)
	newItems.Events = filterEntries(items.Events, cutOff, location)

	return newItems
}
