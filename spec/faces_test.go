package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/entry"
	"time"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Faces Entry", func() {

	var (
		face Face
	)

	BeforeEach( func() {
		face = NewFace(spec.MockClock{}, "aleung")
	})

	Describe("creating a new Faces", func() {
		It("should default the date to today", func() {
			Expect(face.Time).To(Equal(time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)))
			Expect(face.Name).To(BeEmpty())
			Expect(face.Id).To(BeEmpty())
		})
	})

	Describe("creating a new Faces", func() {
		Context("validating when not all mandatory fields are set", func() {
			It("should return false", func() {
				Expect(face.Validate()).To(BeFalse())
			})
		})

		Context("validating when all mandatory fields are set", func() {
			It("should return true", func() {
				face.Name = "some name"
				Expect(face.Validate()).To(BeTrue())
			})
		})
	})

	Describe("when printing out a face", func() {
		It("should print the face", func() {
			face.Name = "some name"
			Expect(face.String()).To(Equal("faces\n  *name: some name\n  date: 2015-01-02"))
		})
	})
})
