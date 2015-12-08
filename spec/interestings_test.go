package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/entry"
	"time"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Interestings Entry", func() {

	var (
		interesting Interesting
	)

	BeforeEach( func() {
		interesting = NewInteresting(spec.MockClock{}, "aleung")
	})

	Describe("creating a new Interestings", func() {
		It("should default the date to today", func() {
			Expect(interesting.Time).To(Equal(time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)))
			Expect(interesting.Title).To(BeEmpty())
			Expect(interesting.Body).To(BeEmpty())
			Expect(interesting.Id).To(BeEmpty())
		})
	})

	Describe("creating a new Interestings", func() {
		Context("validating when not all mandatory fields are set", func() {
			It("should return false", func() {
				Expect(interesting.Validate()).To(BeFalse())
			})
		})

		Context("validating when all mandatory fields are set", func() {
			It("should return true", func() {
				interesting.Title = "some title"
				Expect(interesting.Validate()).To(BeTrue())
			})
		})
	})

	Describe("when printing out an interesting", func() {
		It("should print the interesting", func() {
			interesting.Title = "some title"
			interesting.Body = "some body"
			Expect(interesting.String()).To(Equal("interestings\n  *title: some title\n  body: some body\n  date: 2015-01-02"))
		})
	})
})
