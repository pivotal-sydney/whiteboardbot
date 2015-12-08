package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/entry"
	"time"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Entry", func() {

	var (
		face Face
		clock spec.MockClock
		entry Entry
	)

	BeforeEach( func() {
		clock = spec.MockClock{}
		entry = Entry{}
		face = NewFace(clock, "aleung")
	})

	Context("Entry", func() {
		Describe("validating when not all mandatory fields are set", func() {
			It("should return false", func() {
				Expect(entry.Validate()).To(BeFalse())
			})
		})

		Describe("validating when all mandatory fields are set", func() {
			It("should return true", func() {
				entry.Title = "some name"
				Expect(entry.Validate()).To(BeTrue())
			})
		})
	})

	Context("Faces", func() {
		Describe("creating a new Face", func() {
			It("should have proper defaults", func() {
				Expect(face.Time).To(Equal(time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)))
				Expect(face.Title).To(BeEmpty())
				Expect(face.Body).To(BeEmpty())
				Expect(face.Author).To(Equal("aleung"))
				Expect(face.Id).To(BeEmpty())
			})
		})

		Describe("when printing out a face", func() {
			It("should print the face", func() {
				face.Title = "some name"
				Expect(face.String()).To(Equal("faces\n  *name: some name\n  date: 2015-01-02"))
			})
		})
	})
})
