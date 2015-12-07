package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xtreme-andleung/whiteboardbot/entry"
	. "github.com/xtreme-andleung/whiteboardbot/rest"
	"github.com/xtreme-andleung/whiteboardbot/spec"
	"os"
)

var _ = Describe("Whiteboard Request", func() {

	var (
		face *entry.Face
	)

	BeforeEach(func() {
		face = entry.NewFace(spec.MockClock{}, "aleung")
		face.Name = "Dariusz"
		os.Setenv("token", "token")
	})

	Describe("when creating a NewCreateFaceRequest", func() {
		It("should create request", func() {
			var request = NewCreateFaceRequest(face)
			Expect(request.Utf8).To(Equal(""))
			Expect(request.Method).To(Equal(""))
			Expect(request.Token).To(Equal("token"))
			Expect(request.Commit).To(Equal("Create New Face"))
			Expect(request.Id).To(Equal(""))
			Expect(request.Item.StandupId).To(Equal(1))
			Expect(request.Item.Title).To(Equal("Dariusz"))
			Expect(request.Item.Date).To(Equal("2015-01-02"))
			Expect(request.Item.PostId).To(Equal(""))
			Expect(request.Item.Public).To(Equal(false))
			Expect(request.Item.Kind).To(Equal("New Face"))
			Expect(request.Item.Description).To(Equal(""))
			Expect(request.Item.Author).To(Equal("aleung"))
		})
	})
})

