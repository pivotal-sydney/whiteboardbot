package spec_test

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"github.com/xtreme-andleung/whiteboardbot/model"
)

var _ = Describe("Faces Integration", func() {
	var (
		slackClient spec.MockSlackClient
		clock       spec.MockClock
		restClient  spec.MockRestClient
		whiteboard  WhiteboardApp

		registrationEvent MessageEvent
		newFaceEvent      MessageEvent
		setNameEvent      MessageEvent
		setDateEvent      MessageEvent
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		whiteboard = WhiteboardApp{SlackClient: &slackClient, Clock: clock, RestClient: &restClient, Store: &spec.MockStore{}, EntryMap: make(map[string]model.EntryType)}

		registrationEvent = CreateMessageEvent("wb r 1")
		newFaceEvent = CreateMessageEvent("wb faces")
		setNameEvent = CreateMessageEvent("wb name Dariusz Lorenc")
		setDateEvent = CreateMessageEvent("wb date 2015-12-01")

		whiteboard.ParseMessageEvent(&registrationEvent)
	})

	Describe("with faces keyword", func() {
		It("should begin creating a new face entry and respond with face string", func() {
			whiteboard.ParseMessageEvent(&newFaceEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Face{}))
			Expect(slackClient.Status).To(BeEmpty())
		})
	})

	Context("setting a name detail", func() {
		Describe("with a new face entry started", func() {
			BeforeEach(func() {
				whiteboard.ParseMessageEvent(&newFaceEvent)
			})
			Describe("with correct keyword", func() {
				It("should set the name of the entry and respond with face string", func() {
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(slackClient.EntryType.GetEntry().Title).To(Equal("Dariusz Lorenc"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP))
				})
				It("should post new face entry to whiteboard since all mandatory fields are set", func() {
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					Expect(restClient.Request.Commit).To(Equal("Create New Face"))
				})
				It("should update existing face entry in the whiteboard ", func() {
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					setNameEvent.Text = "wb name updated name"
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
					Expect(restClient.Request.Method).To(Equal("patch"))
					Expect(restClient.Request.Commit).To(Equal("Update New Face"))
					Expect(restClient.Request.Item.Title).To(Equal("updated name"))
					Expect(restClient.Request.Id).To(Equal("1"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP))
				})
				It("should not update existing face entry in the whiteboard when incorrect keyword", func() {
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
					setNameEvent.Text = "wb invalid"
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(1))
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
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(slackClient.Message).To(Equal("Face does not have a body! You idiot."))
				})
			})
		})
	})
})
