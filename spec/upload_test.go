package spec

import (
	. "github.com/benjamintanweihao/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

var _ = Describe("Upload Integration", func() {
	var (
		whiteboard WhiteboardApp
		slackClient *MockSlackClient
		uploadEvent MessageEvent
		file *File
	)

	BeforeEach(func() {
		whiteboard = createWhiteboardAndRegisterStandup(1)
		slackClient = whiteboard.SlackClient.(*MockSlackClient)

		file = &File{}
		file.Permalink = "http://upload/link"
		file.InitialComment = Comment{Comment: "Body of the event"}
		file.Title = "wb i My Title"
		uploadEvent = MessageEvent{Msg: Msg{Upload: true, File: file, Channel: "whiteboard-sydney"}}
	})

	Describe("when uploading an image", func() {
		It("should create an entry using the title command and set the body to the comment with file URL", func() {
			whiteboard.ParseMessageEvent(&uploadEvent)
			Expect(slackClient.Entry.ItemKind).To(Equal("Interesting"))
			Expect(slackClient.Entry.Title).To(Equal("My Title"))
			Expect(slackClient.Entry.Body).To(Equal("Body of the event\n<img src=\"http://upload/link\" style=\"max-width: 500px\">"))
			Expect(slackClient.Status).To(Equal(THUMBS_UP + "_Now go update the details. Need help?_ `wb ?`\n\nINTERESTING\n"))
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
