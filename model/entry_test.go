package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"github.com/pivotal-sydney/whiteboardbot/spec"
	"os"
)

var _ = Describe("Entry", func() {

	var (
		clock spec.MockClock
		entry *Entry
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		entry = NewEntry(clock, "aleung", "title", Standup{Id: 1, TimeZone: "Australia/Sydney"}, "Event")
		os.Setenv("WB_AUTH_TOKEN", "token")
	})

	Context("creating a new Entry", func() {
		It("should have proper defaults", func() {
			Expect(entry.Date).To(Equal("2015-01-02"))
			Expect(entry.Title).To(Equal("title"))
			Expect(entry.Body).To(BeEmpty())
			Expect(entry.Author).To(Equal("aleung"))
			Expect(entry.Id).To(BeEmpty())
			Expect(entry.StandupId).To(Equal(1))
		})
		Describe("with different time zone", func() {
			It("should use the correct time zone for date", func() {
				entry = NewEntry(clock, "aleung", "title", Standup{Id: 1, TimeZone: "America/New_York"}, "Event")
				Expect(entry.Date).To(Equal("2015-01-01"))
			})
		})
	})

	Describe("validating when not all mandatory fields are set", func() {
		It("should return false", func() {
			entry.Title = ""
			Expect(entry.Validate()).To(BeFalse())
		})
	})

	Describe("validating when all mandatory fields are set", func() {
		It("should return true", func() {
			Expect(entry.Validate()).To(BeTrue())
		})
	})
})
