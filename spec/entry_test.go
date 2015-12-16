package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"os"
	"time"
)

var _ = Describe("Entry", func() {

	var (
		clock spec.MockClock
		entry *Entry
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		entry = NewEntry(clock, "aleung", "title", 1)
		os.Setenv("WB_AUTH_TOKEN", "token")
	})

	Describe("creating a new Entry", func() {
		It("should have proper defaults", func() {
			Expect(entry.Date).To(Equal(time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)))
			Expect(entry.Title).To(Equal("title"))
			Expect(entry.Body).To(BeEmpty())
			Expect(entry.Author).To(Equal("aleung"))
			Expect(entry.Id).To(BeEmpty())
			Expect(entry.StandupId).To(Equal(1))
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
