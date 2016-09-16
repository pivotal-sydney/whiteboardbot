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
	var (
		token          string
		writer         *httptest.ResponseRecorder
		mockWhiteBoard MockQuietWhiteboard
		handlerFunc    http.HandlerFunc
	)

	BeforeEach(func() {
		token = os.Getenv("SLACK_TOKEN")

		os.Setenv("SLACK_TOKEN", "123")

		store := MockStore{}
		writer = httptest.NewRecorder()
		whiteboardServer := WhiteboardHttpServer{Store: &store}

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

		AssertDoesNotInvokeHandleInput := func(params *map[string]string) func() {
			return func() {
				request := makeRequest(*params)
				handlerFunc.ServeHTTP(writer, request)

				Expect(mockWhiteBoard.HandleInputCalled).To(BeFalse())
			}
		}

		AssertReturns403Forbidden := func(params *map[string]string) func() {
			return func() {
				request := makeRequest(*params)
				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Code).To(Equal(http.StatusForbidden))
			}
		}

		AssertReturnsErrorMessage := func(params *map[string]string) func() {
			return func() {
				request := makeRequest(*params)
				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Body.String()).To(Equal("Uh-oh, something went wrong... sorry!"))
			}
		}

		Context("when the SLACK_TOKEN environment variable is blank", func() {
			var params map[string]string

			BeforeEach(func() {
				params = map[string]string{"token": ""}
				os.Unsetenv("SLACK_TOKEN")
			})

			It("does not invoke QuietWhiteboard.HandleInput", AssertDoesNotInvokeHandleInput(&params))
			It("returns a 403 Forbidden", AssertReturns403Forbidden(&params))
			It("returns an error message'", AssertReturnsErrorMessage(&params))
		})

		Context("when the token is invalid", func() {
			var params map[string]string

			BeforeEach(func() {
				params = map[string]string{"token": "invalid"}
			})

			It("does not invoke QuietWhiteboard.HandleInput", AssertDoesNotInvokeHandleInput(&params))
			It("returns a 403 Forbidden", AssertReturns403Forbidden(&params))
			It("returns an error message'", AssertReturnsErrorMessage(&params))
		})
	})
})
