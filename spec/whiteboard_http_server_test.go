package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/http"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("WhiteboardHttpServer", func() {
	Describe("HandleRequest", func() {
		It("invokes QuietWhiteboard.HandleInput with payload text", func() {
			// ..
		})

		It("returns the JSON representation of the QuietWhiteboard response", func() {

		})

		Context("when serializing the response fails", func() {

		})

		Context("when the token is invalid", func() {
			It("does not invoke QuietWhiteboard.HandleInput", func() {

			})

			It("returns a 403 Forbidden", func() {
				restClient := MockRestClient{}
				store := MockStore{}

				whiteboardServer := NewWhiteboardHttpServer(&store)
				whiteboard := NewQuietWhiteboard(&restClient, &store)

				handlerFunc := whiteboardServer.NewHandleRequest(whiteboard)

				writer := httptest.NewRecorder()

				request, err := http.NewRequest("POST", "/", nil)

				Expect(err).NotTo(HaveOccurred())

				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Code).To(Equal(http.StatusForbidden))
			})
		})
	})
})
