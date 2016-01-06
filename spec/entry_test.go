package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"os"
	"github.com/xtreme-andleung/whiteboardbot/model"
	"fmt"
)

var _ = Describe("Entry", func() {

	var (
		clock spec.MockClock
		entry *Entry
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		entry = NewEntry(clock, "aleung", "title", Standup{Id: 1, TimeZone: "Australia/Sydney"})
		os.Setenv("WB_AUTH_TOKEN", "token")
	})

	Context("creating a new Entry", func() {
		It("should have proper defaults", func() {
			Expect(entry.Date).To(Equal("2015-01-02"))
			Expect(entry.Title).To(Equal("title"))
			Expect(entry.Body).To(BeEmpty())
			Expect(entry.Author).To(Equal("aleung"))
			Expect(entry.Id).To(BeEmpty())
			Expect(entry.StandupId).To(Equal(1))
		})
		Describe("with different time zone", func() {
			It("should use the correct time zone for date", func() {
				entry = NewEntry(clock, "aleung", "title", Standup{Id: 1, TimeZone: "America/New_York"})
				Expect(entry.Date).To(Equal("2015-01-01"))
			})
		})
	})


	Describe("validating when not all mandatory fields are set", func() {
		It("should return false", func() {
			entry.Title = ""
			Expect(entry.Validate()).To(BeFalse())
		})
	})

	Describe("validating when all mandatory fields are set", func() {
		It("should return true", func() {
			Expect(entry.Validate()).To(BeTrue())
		})
	})

	Context("StandupItems string methods", func() {
		var (
			items StandupItems
		)
		BeforeEach(func() {
			items = model.StandupItems{}
			items.Faces = []model.Entry{model.Entry{Title: "Dariusz", Date: "2015-12-03", Author: "Andrew"}, model.Entry{Title: "Andrew", Date: "2015-12-03", Author: "Dariusz"}}
			items.Interestings = []model.Entry{model.Entry{Title: "Something interesting", Body: "link", Author: "Mik", Date: "2015-12-03"}, model.Entry{Title: "Something else interesting", Body: "link", Author: "Mik", Date: "2016-12-03"}}
			items.Events = []model.Entry{model.Entry{Title: "Another meetup", Body: "link", Author: "Dariusz", Date: "2015-12-03"}, model.Entry{Title: "Another bloody meetup", Body: "link", Author: "Dariusz", Date: "2025-12-03"}}
			items.Helps = []model.Entry{model.Entry{Title: "Help me!", Author: "Lawrence", Date: "2015-12-03"}, model.Entry{Title: "Help me again!", Author: "Lawrence", Date: "2016-12-03"}}
		})


		Describe("convert standup items to string", func() {
			It("should return string representation of standup items", func() {

				itemsString := items.String()
				Expect(itemsString).To(Equal(fmt.Sprintf("%v\n%v\n%v\n%v", items.FacesString(), items.InterestingsString(), items.HelpsString(), items.EventsString())))
			})
		})

		Describe("convert standup faces items to string", func() {
			It("should print faces in presentation mode", func() {
				itemsString := items.FacesString()
				Expect(itemsString).To(Equal("NEW FACES\n\n" + Face{&items.Faces[0]}.String() + "\n\n" + Face{&items.Faces[1]}.String()))
			})
		})

		Describe("convert standup interestings items to string", func() {
			It("should print interestings in presentation mode", func() {
				itemsString := items.InterestingsString()
				Expect(itemsString).To(Equal("INTERESTINGS\n\n" + Interesting{&items.Interestings[0]}.String() + "\n\n" + Interesting{&items.Interestings[1]}.String()))
			})
		})

		Describe("convert standup helps items to string", func() {
			It("should print helps in presentation mode", func() {
				itemsString := items.HelpsString()
				Expect(itemsString).To(Equal("HELPS\n\n" + Help{&items.Helps[0]}.String() + "\n\n" + Help{&items.Helps[1]}.String()))
			})
		})

		Describe("convert standup events items to string", func() {
			It("should print events in presentation mode", func() {
				itemsString := items.EventsString()
				Expect(itemsString).To(Equal("EVENTS\n\n" + Event{&items.Events[0]}.String() + "\n\n" + Event{&items.Events[1]}.String()))
			})
		})
	})
})
