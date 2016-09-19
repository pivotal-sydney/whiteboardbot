package model_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-sydney/whiteboardbot/model"
	. "github.com/pivotal-sydney/whiteboardbot/model"
)

var _ = Describe("StandupItems", func() {

	var items StandupItems

	BeforeEach(func() {
		items = model.StandupItems{}
		items.Faces = []model.Entry{{Title: "Dariusz", Date: "2015-12-03", Author: "Andrew"}, {Title: "Andrew", Date: "2015-12-03", Author: "Dariusz"}}
		items.Interestings = []model.Entry{{Title: "Something interesting", Body: "link", Author: "Mik", Date: "2015-12-03"}, {Title: "Something else interesting", Body: "link", Author: "Mik", Date: "2016-12-03"}}
		items.Events = []model.Entry{{Title: "Another meetup", Body: "link", Author: "Dariusz", Date: "2015-12-03"}, {Title: "Another bloody meetup", Body: "link", Author: "Dariusz", Date: "2025-12-03"}}
		items.Helps = []model.Entry{{Title: "Help me!", Author: "Lawrence", Date: "2015-12-03"}, {Title: "Help me again!", Author: "Lawrence", Date: "2016-12-03"}}
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
})
