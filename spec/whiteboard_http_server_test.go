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
	"os"
	"strconv"
)

type MockQuietWhiteboard struct {
	HandleInputCalled bool
}

func (mqw *MockQuietWhiteboard) HandleInput(input string) Response {
	mqw.HandleInputCalled = true
	return Response{}
}

func makeRequest(params map[string]string) *http.Request {
	data := url.Values{}

	data.Set("token", "123")
	for k, v := range params {
		data.Set(k, v)
	}

	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	return request
}

var _ = Describe("WhiteboardHttpServer", func() {
	var token string

	var writer *httptest.ResponseRecorder
	var mockWhiteBoard MockQuietWhiteboard
	var handlerFunc http.HandlerFunc

	BeforeEach(func() {
		token = os.Getenv("SLACK_TOKEN")

		os.Setenv("SLACK_TOKEN", "123")

		store := MockStore{}
		writer = httptest.NewRecorder()
		whiteboardServer := NewWhiteboardHttpServer(&store)

		mockWhiteBoard = MockQuietWhiteboard{}
		handlerFunc = whiteboardServer.NewHandleRequest(&mockWhiteBoard)
	})

	AfterEach(func() {
		os.Setenv("SLACK_TOKEN", token)
	}, 0)

	Describe("HandleRequest", func() {
		It("invokes QuietWhiteboard.HandleInput with payload text", func() {
			params := make(map[string]string)
			request := makeRequest(params)

			handlerFunc.ServeHTTP(writer, request)

			Expect(mockWhiteBoard.HandleInputCalled).To(BeTrue())
		})

		It("returns the JSON representation of the QuietWhiteboard response", func() {
			params := make(map[string]string)
			request := makeRequest(params)

			handlerFunc.ServeHTTP(writer, request)

			Expect(writer.Body.String()).To(Equal(`{"text":""}`))
		})

		Context("when the SLACK_TOKEN environment variable is blank", func() {
			It("does not invoke QuietWhiteboard.HandleInput", func() {
				os.Unsetenv("SLACK_TOKEN")
				params := map[string]string{"token": ""}
				request := makeRequest(params)

				handlerFunc.ServeHTTP(writer, request)

				Expect(mockWhiteBoard.HandleInputCalled).To(BeFalse())
			})

			It("returns a 403 Forbidden", func() {
				os.Unsetenv("SLACK_TOKEN")
				params := map[string]string{"token": ""}
				request := makeRequest(params)

				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("when the token is invalid", func() {
			It("does not invoke QuietWhiteboard.HandleInput", func() {
				params := map[string]string{"token": "invalid"}
				request := makeRequest(params)

				handlerFunc.ServeHTTP(writer, request)

				Expect(mockWhiteBoard.HandleInputCalled).To(BeFalse())
			})

			It("returns a 403 Forbidden", func() {
				params := map[string]string{"token": "invalid"}
				request := makeRequest(params)

				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Code).To(Equal(http.StatusForbidden))
			})
		})
	})
})
