package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/nlopes/slack"
)
var _ = Describe("Slack Client", func() {

	var slackUser *slack.User

	Context("when user profile is not available", func() {
		It("should use user name as author", func() {
			slackUser = &slack.User{Name: "aleung"}
			author := app.GetAuthor(slackUser)
			Expect(author).To(Equal("aleung"))
		})
	})
	Context("when user profile is available", func() {
		It("should use user real name as author", func() {
			slackUser = &slack.User{Name: "aleung", Profile: slack.UserProfile{RealName: "Andrew Leung"}}
			author := app.GetAuthor(slackUser)
			Expect(author).To(Equal("Andrew Leung"))
		})
	})
})



