package model_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"github.com/pivotal-sydney/whiteboardbot/spec"
)

var _ = Describe("StandupItems", func() {

	var (
		items StandupItems
		clock spec.MockClock
	)

	BeforeEach(func() {
		items = StandupItems{}

		clock = spec.MockClock{}
		today := clock.Now()
		tomorrow := today.AddDate(0, 0, 1)
		dayAfterTomorrow := tomorrow.AddDate(0, 0, 1)

		todayStr := today.Format("2006-01-02")
		tomorrowStr := tomorrow.Format("2006-01-02")
		dayAfterTomorrowStr := dayAfterTomorrow.Format("2006-01-02")

		items.Faces = []Entry{
			{Title: "Dariusz", Date: todayStr, Author: "Andrew"},
			{Title: "Andrew", Date: tomorrowStr, Author: "Dariusz"},
		}

		items.Interestings = []Entry{
			{Title: "Something interesting", Body: "link", Author: "Mik", Date: todayStr},
			{Title: "Something else interesting", Body: "link", Author: "Mik", Date: "Dinner on Friday?"},
		}

		items.Events = []Entry{
			{Title: "Another meetup", Body: "link", Author: "Dariusz", Date: tomorrowStr},
			{Title: "Another bloody meetup", Body: "link", Author: "Dariusz", Date: dayAfterTomorrowStr},
		}

		items.Helps = []Entry{
			{Title: "Help me!", Author: "Lawrence", Date: todayStr},
			{Title: "Help me again!", Author: "Lawrence", Date: todayStr},
		}
	})

	Describe("convert standup items to string", func() {
		It("should return string representation of standup items", func() {
			itemsString := items.String()
			Expect(itemsString).To(Equal(fmt.Sprintf(">>>— — —\n \n \n \n%v\n \n \n \n%v\n \n \n \n%v\n \n \n \n%v\n \n \n \n— — —\n:clap:", items.FacesString(), items.HelpsString(), items.InterestingsString(), items.EventsString())))
		})
	})

	Describe("convert standup faces items to string", func() {
		It("should print faces in presentation mode", func() {
			itemsString := items.FacesString()
			Expect(itemsString).To(Equal("NEW FACES\n\n" + Face{&items.Faces[0]}.String() + "\n \n" + Face{&items.Faces[1]}.String()))
		})
	})

	Describe("convert standup interestings items to string", func() {
		It("should print interestings in presentation mode", func() {
			itemsString := items.InterestingsString()
			Expect(itemsString).To(Equal("INTERESTINGS\n\n" + Interesting{&items.Interestings[0]}.String() + "\n \n" + Interesting{&items.Interestings[1]}.String()))
		})
	})

	Describe("convert standup helps items to string", func() {
		It("should print helps in presentation mode", func() {
			itemsString := items.HelpsString()
			Expect(itemsString).To(Equal("HELPS\n\n" + Help{&items.Helps[0]}.String() + "\n \n" + Help{&items.Helps[1]}.String()))
		})
	})

	Describe("convert standup events items to string", func() {
		It("should print events in presentation mode", func() {
			itemsString := items.EventsString()
			Expect(itemsString).To(Equal("EVENTS\n\n" + Event{&items.Events[0]}.String() + "\n \n" + Event{&items.Events[1]}.String()))
		})
	})

	Describe("Filter", func() {
		It("lists all items", func() {
			todaysItems := items.Filter(1, clock, "Australia/Sydney")

			Expect(todaysItems.Faces).To(Equal(items.Faces[:1]))
			Expect(todaysItems.Helps).To(Equal(items.Helps))
			Expect(todaysItems.Interestings).To(Equal(items.Interestings))
			Expect(todaysItems.Events).To(BeEmpty())

			tomorrowsItems := items.Filter(2, clock, "Australia/Sydney")

			Expect(tomorrowsItems.Faces).To(Equal(items.Faces))
			Expect(tomorrowsItems.Helps).To(Equal(items.Helps))
			Expect(tomorrowsItems.Interestings).To(Equal(items.Interestings))
			Expect(tomorrowsItems.Events).To(Equal(items.Events[:1]))

			allItems := items.Filter(3, clock, "Australia/Sydney")

			Expect(allItems.Faces).To(Equal(items.Faces))
			Expect(allItems.Helps).To(Equal(items.Helps))
			Expect(allItems.Interestings).To(Equal(items.Interestings))
			Expect(allItems.Events).To(Equal(items.Events))
		})

		Context("when the user's time zone is invalid", func() {
			It("uses the system time zone", func() {
				newItems := items.Filter(1, clock, "bork bork bork")

				Expect(newItems.Faces).To(Equal(items.Faces[:1]))
				Expect(newItems.Helps).To(Equal(items.Helps))
				Expect(newItems.Interestings).To(Equal(items.Interestings))
				Expect(newItems.Events).To(BeEmpty())
			})
		})
	})
})
