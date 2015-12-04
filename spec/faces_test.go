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
		face = NewFace(spec.MockClock{})
	})

	Describe("Creating a new Faces", func() {
		It("Should default the date to today", func() {
			Expect(face.GetDate()).To(Equal(time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)))
			Expect(face.GetName()).To(BeEmpty())
		})
	})

	Describe("Creating a new Faces", func() {
		Context("When setting a date", func() {
			It("Should update the date", func() {
				now := time.Now()
				face.SetDate(now)
				Expect(face.GetDate()).To(Equal(now))
			})
		})
		Context("When setting a name", func() {
			It("Should update the name", func() {
				name := "new name"
				face.SetName(name)
				Expect(face.GetName()).To(Equal(name))
			})
		})

		Context("Validating when not all mandatory fields are set", func() {
			It("It should return false", func() {
				Expect(face.Validate()).To(BeFalse())
			})
		})

		Context("Validating when all mandatory fields are set", func() {
			It("It should return true", func() {
				face.SetName("some name")
				Expect(face.Validate()).To(BeTrue())
			})
		})
	})

	Describe("When printing out a face", func() {
		It("Should print the face", func() {
			face.SetName("some name")
			Expect(face.String()).To(Equal("faces\n*name: some name\ndate: 2015-01-02"))
		})
	})

})
