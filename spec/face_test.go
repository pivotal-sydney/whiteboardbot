package spec

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

var _ = Describe("Faces Integration", func() {
	var (
		whiteboard WhiteboardApp
		slackClient *MockSlackClient
		restClient *MockRestClient

		newFaceEvent, newFaceWithTitleEvent, setNameEvent, setDateEvent MessageEvent
	)

	BeforeEach(func() {
		whiteboard = createWhiteboardAndRegisterStandup(1)
		slackClient = whiteboard.SlackClient.(*MockSlackClient)
		restClient = whiteboard.RestClient.(*MockRestClient)

		newFaceEvent = createMessageEvent("wb faces")
		newFaceWithTitleEvent = createMessageEvent("wb faces Andrew Leung")
		setNameEvent = createMessageEvent("wb name Dariusz Lorenc")
		setDateEvent = createMessageEvent("wb date 2015-12-01")
	})

	Describe("with faces keyword without title", func() {
		It("should respond missing title", func() {
			whiteboard.ParseMessageEvent(&newFaceEvent)
			Expect(slackClient.Message).To(Equal("Hey, next time add a title along with your entry!\nLike this: `wb i My title`\nNeed help? try `wb ?`"))
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
					Expect(slackClient.Entry.Title).To(Equal("Andrew Leung"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP + "_Now go update the details. Need help?_ `wb ?`\n\nNEW FACE\n"))
					Expect(restClient.PostCalledCount).To(Equal(1))
					Expect(restClient.Request.Commit).To(Equal("Create New Face"))
				})

				It("should update existing face entry in the whiteboard ", func() {
					whiteboard.ParseMessageEvent(&setNameEvent)
					Expect(restClient.PostCalledCount).To(Equal(2))
					Expect(slackClient.Entry.Title).To(Equal("Dariusz Lorenc"))
					Expect(restClient.Request.Method).To(Equal("patch"))
					Expect(restClient.Request.Commit).To(Equal("Update New Face"))
					Expect(restClient.Request.Item.Title).To(Equal("Dariusz Lorenc"))
					Expect(restClient.Request.Id).To(Equal("1"))
					Expect(slackClient.Status).To(Equal(THUMBS_UP + "NEW FACE\n"))
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