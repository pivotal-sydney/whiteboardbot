package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	. "github.com/nlopes/slack"
	 ."github.com/xtreme-andleung/whiteboardbot/app"
)

var _ = Describe("Upload Integration", func() {

	var (
		slackClient spec.MockSlackClient
		clock spec.MockClock
		restClient spec.MockRestClient
		uploadEvent MessageEvent
		file *File
	)

	BeforeEach(func() {
		slackClient = spec.MockSlackClient{}
		clock = spec.MockClock{}
		restClient = spec.MockRestClient{}
		file = &File{}
		file.URL = "http://upload/link"
		file.InitialComment = Comment{Comment: "Body of the event"}
		file.Title = "wb i My Title"
		uploadEvent = MessageEvent{Msg: Msg{Upload: true, File: file}}
	})

	Describe("when uploading an image", func() {
		It("should create an entry using the title command and set the body to the comment with file URL", func() {
			ParseMessageEvent(&slackClient, &restClient, clock, &uploadEvent)
			Expect(slackClient.Message).To(Equal("interestings\n  *title: My Title\n  body: Body of the event\n![](http://upload/link)\n  date: 2015-01-02\nitem created"))
		})
		Context("with invalid keyword", func() {
			BeforeEach(func() {
				file.Title = "wb nonKeyword"
			})

			It("should handle default response", func() {
				ParseMessageEvent(&slackClient, &restClient, clock, &uploadEvent)
				Expect(slackClient.Message).To(Equal("aleung no you nonKeyword"))
			})
		})
	})
})