package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/model"
	"time"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"os"
)

var _ = Describe("Entry", func() {

	var (
		face Face
		interesting Interesting
		clock spec.MockClock
		entry Entry
	)

	BeforeEach(func() {
		clock = spec.MockClock{}
		entry = Entry{}
		face = NewFace(clock, "aleung")
		interesting = NewInteresting(clock, "aleung")
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

		Context("when making requets", func() {

			BeforeEach(func() {
				face.Title = "Dariusz"
				face.Id = "123"
				os.Setenv("WB_AUTH_TOKEN", "token")
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

	Context("Interestings", func() {
		Describe("creating a new Interesting", func() {
			It("should have proper defaults", func() {
				Expect(interesting.Time).To(Equal(time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)))
				Expect(interesting.Title).To(BeEmpty())
				Expect(interesting.Body).To(BeEmpty())
				Expect(interesting.Author).To(Equal("aleung"))
				Expect(interesting.Id).To(BeEmpty())
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
})
