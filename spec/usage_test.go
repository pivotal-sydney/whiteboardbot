package spec

import (
	. "github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
)

var _ = Describe("Usage Integration", func() {
	var (
		whiteboard  WhiteboardApp
		slackClient *MockSlackClient
		usageEvent  MessageEvent
	)

	BeforeEach(func() {
		whiteboard = createWhiteboardAndRegisterStandup(1)
		slackClient = whiteboard.SlackClient.(*MockSlackClient)
		usageEvent = createMessageEvent("wb ?")
	})

	Describe("when question mark command is send", func() {
		It("should respond with usage screen", func() {
			whiteboard.ParseMessageEvent(&usageEvent)
			Expect(slackClient.Message).Should(Equal(USAGE))
		})
	})
})
