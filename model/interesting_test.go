package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"github.com/pivotal-sydney/whiteboardbot/spec"
	"os"
)

var _ = Describe("Interesting", func() {

	var (
		interesting Interesting
		clock       spec.MockClock
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		interesting = NewInteresting(clock, "aleung", "title", Standup{Id: 1, TimeZone: "Australia/Sydney"}).(Interesting)
		os.Setenv("WB_AUTH_TOKEN", "token")
	})

	Describe("creating a new Interesting", func() {
		It("should have proper defaults", func() {
			Expect(interesting.Date).To(Equal("2015-01-02"))
			Expect(interesting.Title).To(Equal("title"))
			Expect(interesting.Body).To(BeEmpty())
			Expect(interesting.Author).To(Equal("aleung"))
			Expect(interesting.Id).To(BeEmpty())
			Expect(interesting.StandupId).To(Equal(1))
		})
	})

	Describe("when printing out an interesting", func() {
		It("should print the interesting", func() {
			interesting.Title = "some title"
			interesting.Body = "some body"
			interesting.Author = "some author"
			Expect(interesting.String()).To(Equal("*some title*\nsome body\n[some author]\n02 Jan 2015"))
		})
	})

	Context("when making requets", func() {
		BeforeEach(func() {
			interesting.Title = "Dariusz"
			interesting.Id = "123"
			interesting.Body = "Body Text"
		})

		Describe("create request", func() {
			It("should populate request with fields", func() {
				request := interesting.MakeCreateRequest()
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
				Expect(request.Item.Kind).To(Equal("Interesting"))
				Expect(request.Item.Description).To(Equal("Body Text"))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})

		Describe("update request", func() {
			It("should populate request with fields", func() {
				request := interesting.MakeUpdateRequest()
				Expect(request.Utf8).To(Equal(""))
				Expect(request.Method).To(Equal("patch"))
				Expect(request.Token).To(Equal("token"))
				Expect(request.Commit).To(Equal("Update Item"))
				Expect(request.Id).To(Equal(interesting.Id))
				Expect(request.Item.StandupId).To(Equal(1))
				Expect(request.Item.Title).To(Equal("Dariusz"))
				Expect(request.Item.Date).To(Equal("2015-01-02"))
				Expect(request.Item.PostId).To(Equal(""))
				Expect(request.Item.Public).To(Equal("false"))
				Expect(request.Item.Kind).To(Equal("Interesting"))
				Expect(request.Item.Description).To(Equal("Body Text"))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})
	})
})
