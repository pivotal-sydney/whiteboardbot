package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"github.com/pivotal-sydney/whiteboardbot/spec"
	"os"
)

var _ = Describe("Event", func() {

	var (
		event Event
		clock spec.MockClock
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		event = NewEvent(clock, "aleung", "title", Standup{Id: 1, TimeZone: "Australia/Sydney"})
		os.Setenv("WB_AUTH_TOKEN", "token")
	})

	Describe("creating a new Event", func() {
		It("should have proper defaults", func() {
			Expect(event.Date).To(Equal("2015-01-02"))
			Expect(event.Title).To(Equal("title"))
			Expect(event.Body).To(BeEmpty())
			Expect(event.Author).To(Equal("aleung"))
			Expect(event.Id).To(BeEmpty())
			Expect(event.StandupId).To(Equal(1))
		})
	})

	Describe("when printing out an event", func() {
		It("should print the event", func() {
			event.Title = "some title"
			event.Body = "some body"
			Expect(event.String()).To(Equal("events\n  *title: some title\n  body: some body\n  date: 2015-01-02"))
		})
	})

	Context("when making requets", func() {
		BeforeEach(func() {
			event.Title = "Dariusz"
			event.Id = "123"
			event.Body = "Body Text"
		})

		Describe("create request", func() {
			It("should populate request with fields", func() {
				request := event.MakeCreateRequest()
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
				Expect(request.Item.Kind).To(Equal("Event"))
				Expect(request.Item.Description).To(Equal("Body Text"))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})

		Describe("update request", func() {
			It("should populate request with fields", func() {
				request := event.MakeUpdateRequest()
				Expect(request.Utf8).To(Equal(""))
				Expect(request.Method).To(Equal("patch"))
				Expect(request.Token).To(Equal("token"))
				Expect(request.Commit).To(Equal("Update Item"))
				Expect(request.Id).To(Equal(event.Id))
				Expect(request.Item.StandupId).To(Equal(1))
				Expect(request.Item.Title).To(Equal("Dariusz"))
				Expect(request.Item.Date).To(Equal("2015-01-02"))
				Expect(request.Item.PostId).To(Equal(""))
				Expect(request.Item.Public).To(Equal("false"))
				Expect(request.Item.Kind).To(Equal("Event"))
				Expect(request.Item.Description).To(Equal("Body Text"))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})
	})
})
