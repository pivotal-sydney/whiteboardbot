package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"os"
)

var _ = Describe("Help", func() {

	var (
		help Help
		clock spec.MockClock
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		help = NewHelp(clock, "aleung", "title", Standup{Id: 1, TimeZone: "Australia/Sydney"})
		os.Setenv("WB_AUTH_TOKEN", "token")
	})

	Describe("creating a new Help", func() {
		It("should have proper defaults", func() {
			Expect(help.Date).To(Equal("2015-01-02"))
			Expect(help.Title).To(Equal("title"))
			Expect(help.Body).To(BeEmpty())
			Expect(help.Author).To(Equal("aleung"))
			Expect(help.Id).To(BeEmpty())
			Expect(help.StandupId).To(Equal(1))
		})
	})

	Describe("when printing out an help", func() {
		It("should print the help", func() {
			help.Title = "some title"
			help.Body = "some body"
			Expect(help.String()).To(Equal("helps\n  *title: some title\n  body: some body\n  date: 2015-01-02"))
		})
	})

	Context("when making requets", func() {
		BeforeEach(func() {
			help.Title = "Dariusz"
			help.Id = "123"
			help.Body = "Body Text"
		})

		Describe("create request", func() {
			It("should populate request with fields", func() {
				request := help.MakeCreateRequest()
				Expect(request.Utf8).To(Equal(""))
				Expect(request.Method).To(Equal(""))
				Expect(request.Token).To(Equal("token"))
				Expect(request.Commit).To(Equal("Create Item"))
				Expect(request.Id).To(Equal(""))
				Expect(request.Item.StandupId).To(Equal(1))
				Expect(request.Item.Title).To(Equal("Dariusz"))
				Expect(request.Item.Date).To(Equal("2015-01-02"))
				Expect(request.Item.PostId).To(Equal(""))
				Expect(request.Item.Public).To(Equal("false"))
				Expect(request.Item.Kind).To(Equal("Help"))
				Expect(request.Item.Description).To(Equal("Body Text"))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})

		Describe("update request", func() {
			It("should populate request with fields", func() {
				request := help.MakeUpdateRequest()
				Expect(request.Utf8).To(Equal(""))
				Expect(request.Method).To(Equal("patch"))
				Expect(request.Token).To(Equal("token"))
				Expect(request.Commit).To(Equal("Update Item"))
				Expect(request.Id).To(Equal(help.Id))
				Expect(request.Item.StandupId).To(Equal(1))
				Expect(request.Item.Title).To(Equal("Dariusz"))
				Expect(request.Item.Date).To(Equal("2015-01-02"))
				Expect(request.Item.PostId).To(Equal(""))
				Expect(request.Item.Public).To(Equal("false"))
				Expect(request.Item.Kind).To(Equal("Help"))
				Expect(request.Item.Description).To(Equal("Body Text"))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})
	})
})
