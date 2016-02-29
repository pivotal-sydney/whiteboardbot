package spec_test

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	"github.com/pivotal-sydney/whiteboardbot/spec"
	"github.com/pivotal-sydney/whiteboardbot/model"
)

var _ = Describe("Upload Integration", func() {

	var (
		slackClient spec.MockSlackClient
		clock       spec.MockClock
		restClient  spec.MockRestClient
		whiteboard  WhiteboardApp

		uploadEvent       MessageEvent
		registrationEvent MessageEvent
		file              *File
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		whiteboard = WhiteboardApp{SlackClient: &slackClient, Clock: clock, RestClient: &restClient, Store: &spec.MockStore{}, EntryMap: make(map[string]model.EntryType)}

		file = &File{}
		file.URL = "http://upload/link"
		file.InitialComment = Comment{Comment: "Body of the event"}
		file.Title = "wb i My Title"
		uploadEvent = MessageEvent{Msg: Msg{Upload: true, File: file, Channel: "whiteboard-sydney"}}
		registrationEvent = CreateMessageEvent("wb r 1")

		whiteboard.ParseMessageEvent(&registrationEvent)
	})

	Describe("when uploading an image", func() {
		It("should create an entry using the title command and set the body to the comment with file URL", func() {
			whiteboard.ParseMessageEvent(&uploadEvent)
			Expect(slackClient.EntryType).To(BeAssignableToTypeOf(model.Interesting{}))
			Expect(slackClient.EntryType.GetEntry().Title).To(Equal("My Title"))
			Expect(slackClient.EntryType.GetEntry().Body).To(Equal("Body of the event\n<img src=\"http://upload/link\" style=\"max-width: 500px\">"))
			Expect(slackClient.Status).To(Equal("\nitem created"))
		})
		Context("with invalid keyword", func() {
			BeforeEach(func() {
				file.Title = "wb nonKeyword"
			})

			It("should handle default response", func() {
				whiteboard.ParseMessageEvent(&uploadEvent)
				Expect(slackClient.Message).To(Equal("aleung no you nonKeyword"))
			})
		})
	})
})
