package spec_test

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Faces Integration", func() {
	var (
		slackClient spec.MockSlackClient
		clock spec.MockClock
		restClient spec.MockRestClient
		whiteboard WhiteboardApp

		registrationEvent MessageEvent
		newFaceEvent MessageEvent
		newFaceWithTitleEvent MessageEvent
		setNameEvent MessageEvent
		setDateEvent MessageEvent
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		whiteboard = NewWhiteboard(&slackClient, &restClient, clock, &spec.MockStore{})

		registrationEvent = CreateMessageEvent("wb r 1")
		newFaceEvent = CreateMessageEvent("wb faces")
		newFaceWithTitleEvent = CreateMessageEvent("wb faces Andrew Leung")
		setNameEvent = CreateMessageEvent("wb name Dariusz Lorenc")
		setDateEvent = CreateMessageEvent("wb date 2015-12-01")

		whiteboard.ParseMessageEvent(&registrationEvent)
	})

	Describe("with faces keyword without title", func() {
		It("should respond missing title", func() {
			whiteboard.ParseMessageEvent(&newFaceEvent)
			Expect(slackClient.Message).To(Equal("Hey, next time add a title along with your entry!\nLike this: `wb i My title`"))
			Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
		})
	})

	Context("setting a name detail", func() {
		Describe("with a new face entry started", func() {
			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newFaceWithTitleEvent)
			})
			Describe("with correct keyword", func() {
				It("should set the name of the entry and respond with face string", func() {
					Expect(slackClient.EntryType.GetEntry().Title).To(Equal("Andrew Leung"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP))
					Expect(restClient.PostCalledCount).To(Equal(1))
					Expect(restClient.Request.Commit).To(Equal("Create New Face"))
				})
				It("should update existing face entry in the whiteboard ", func() {
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
					Expect(slackClient.EntryType.GetEntry().Title).To(Equal("Dariusz Lorenc"))
					Expect(restClient.Request.Method).To(Equal("patch"))
					Expect(restClient.Request.Commit).To(Equal("Update New Face"))
					Expect(restClient.Request.Item.Title).To(Equal("Dariusz Lorenc"))
					Expect(restClient.Request.Id).To(Equal("1"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP))
				})
				It("should not update existing face entry in the whiteboard when incorrect keyword", func() {
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
					setNameEvent.Text = "wb invalid"
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
				})
			})
			Describe("with incorrect keyword", func() {
				It("should respond with default", func() {
					setNameEvent.Text = "wb nameSomethingWrong"
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(slackClient.Message).To(Equal("aleung no you nameSomethingWrong"))
				})
			})
			Describe("with not allowed keyword", func() {
				It("should respond with random insult", func() {
					setNameEvent.Text = "wb body no body"
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(slackClient.Message).To(Equal("Face does not have a body! Stupid."))
					Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(slackClient.Message).To(Equal("Face does not have a body! You idiot."))
					Expect(slackClient.Status).To(Equal(THUMBS_DOWN))
				})
			})
		})
	})
})
