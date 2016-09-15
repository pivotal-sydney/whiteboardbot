package spec

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/http"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
)

type MockQuietWhiteboard struct {
	HandleInputCalled bool
}

func (mqw MockQuietWhiteboard) HandleInput(input string) Response {
	mqw.HandleInputCalled = true
	return Response{}
}

func makeRequest() *http.Request {
	data := url.Values{}
	data.Add("token", "invalid")

	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	return request
}

var _ = Describe("WhiteboardHttpServer", func() {
	Describe("HandleRequest", func() {
		It("invokes QuietWhiteboard.HandleInput with payload text", func() {
		})

		It("returns the JSON representation of the QuietWhiteboard response", func() {

		})

		Context("when serializing the response fails", func() {

		})

		Context("when the token is invalid", func() {
			It("does not invoke QuietWhiteboard.HandleInput", func() {
				store := MockStore{}
				writer := httptest.NewRecorder()
				whiteboardServer := NewWhiteboardHttpServer(&store)

				mockWhiteBoard := MockQuietWhiteboard{}
				handlerFunc := whiteboardServer.NewHandleRequest(mockWhiteBoard)

				request := makeRequest()

				handlerFunc.ServeHTTP(writer, request)

				Expect(mockWhiteBoard.HandleInputCalled).To(BeFalse())
			})

			It("returns a 403 Forbidden", func() {
				store := MockStore{}
				writer := httptest.NewRecorder()
				whiteboardServer := NewWhiteboardHttpServer(&store)

				mockWhiteBoard := MockQuietWhiteboard{}
				handlerFunc := whiteboardServer.NewHandleRequest(mockWhiteBoard)

				request := makeRequest()

				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Code).To(Equal(http.StatusForbidden))
			})
		})
	})
})
