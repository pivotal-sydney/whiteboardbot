package spec_test

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"github.com/xtreme-andleung/whiteboardbot/model"
)

var _ = Describe("Present Integration", func() {
	var (
		slackClient spec.MockSlackClient
		clock       spec.MockClock
		restClient  spec.MockRestClient
		whiteboard  WhiteboardApp
		presentEvent MessageEvent
		registrationEvent MessageEvent
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		whiteboard = NewWhiteboard(&slackClient, &restClient, clock, &spec.MockStore{})
		registrationEvent = CreateMessageEvent("wb r 1")
		presentEvent = CreateMessageEvent("wb present")
		whiteboard.ParseMessageEvent(&registrationEvent)
	})

	Describe("when present command is sent", func() {
		Context("there is no items in current standup", func() {
			It("should show empty whiteboard", func() {
				whiteboard.ParseMessageEvent(&presentEvent)
				Expect(slackClient.Message).To(Equal("Hey, there's no entries in today's standup yet, why not add some?"))
				Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
			})
		})

		Context("there are items in current standup", func() {
			BeforeEach(func() {
				restClient.StandupItems = model.StandupItems{}
				restClient.StandupItems.Faces = []model.Entry{model.Entry{Title: "Dariusz", Date: "2015-12-03", Author: "Andrew"}}
				restClient.StandupItems.Interestings = []model.Entry{model.Entry{Title: "Something interesting", Body: "link", Author: "Mik", Date: "2015-12-03"}}
				restClient.StandupItems.Events = []model.Entry{model.Entry{Title: "Another meetup", Body: "link", Author: "Dariusz", Date: "2015-12-03"}}
				restClient.StandupItems.Helps = []model.Entry{model.Entry{Title: "Help me!", Author: "Lawrence", Date: "2015-12-03"}}
			})
			It("should display all standup's items", func() {
				whiteboard.ParseMessageEvent(&presentEvent)
				Expect(slackClient.Message).To(Equal(restClient.StandupItems.String()))
			})
		})

	})
})