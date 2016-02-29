package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"github.com/pivotal-sydney/whiteboardbot/spec"
	"os"
	"time"
)

var _ = Describe("Face", func() {

	var (
		face  Face
		clock spec.MockClock
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		face = NewFace(clock, "aleung")
		os.Setenv("WB_AUTH_TOKEN", "token")
	})

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

	Context("when making requets", func() {
		BeforeEach(func() {
			face.Title = "Dariusz"
			face.Id = "123"
		})

		Describe("create request", func() {
			It("should populate request with fields", func() {
				request := face.MakeCreateRequest()
				Expect(request.Utf8).To(Equal(""))
				Expect(request.Method).To(Equal(""))
				Expect(request.Token).To(Equal("token"))
				Expect(request.Commit).To(Equal("Create New Face"))
				Expect(request.Id).To(Equal(""))
				Expect(request.Item.StandupId).To(Equal(1))
				Expect(request.Item.Title).To(Equal("Dariusz"))
				Expect(request.Item.Date).To(Equal("2015-01-02"))
				Expect(request.Item.PostId).To(Equal(""))
				Expect(request.Item.Public).To(Equal("false"))
				Expect(request.Item.Kind).To(Equal("New face"))
				Expect(request.Item.Description).To(Equal(""))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})

		Describe("update request", func() {
			It("should populate request with fields", func() {
				request := face.MakeUpdateRequest()
				Expect(request.Utf8).To(Equal(""))
				Expect(request.Method).To(Equal("patch"))
				Expect(request.Token).To(Equal("token"))
				Expect(request.Commit).To(Equal("Update New Face"))
				Expect(request.Id).To(Equal(face.Id))
				Expect(request.Item.StandupId).To(Equal(1))
				Expect(request.Item.Title).To(Equal("Dariusz"))
				Expect(request.Item.Date).To(Equal("2015-01-02"))
				Expect(request.Item.PostId).To(Equal(""))
				Expect(request.Item.Public).To(Equal("false"))
				Expect(request.Item.Kind).To(Equal("New face"))
				Expect(request.Item.Description).To(Equal(""))
				Expect(request.Item.Author).To(Equal("aleung"))
			})
		})
	})
})
