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

func makeRequest(params map[string]string) *http.Request {
	data := url.Values{}

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
		params         map[string]string
		slackUser      SlackUser
		slackChannel   SlackChannel
	)

	BeforeEach(func() {
		token = os.Getenv("SLACK_TOKEN")

		params = map[string]string{
			"token":        "123",
			"user_name":    "aleung",
			"channel_id":   "C456",
			"channel_name": "sydney-standup",
		}

		slackUser = SlackUser{Username: "aleung", Author: "Andrew Leung", TimeZone: "Australia/Sydney"}
		slackChannel = SlackChannel{ChannelId: "C456", ChannelName: "sydney-standup"}

		os.Setenv("SLACK_TOKEN", "123")

		slackClient := MockSlackClient{}

		writer = httptest.NewRecorder()
		whiteboardServer := WhiteboardHttpServer{SlackClient: &slackClient}

		mockWhiteBoard = MockQuietWhiteboard{}
		handlerFunc = whiteboardServer.NewHandleRequest(&mockWhiteBoard)
	})

	AfterEach(func() {
		os.Setenv("SLACK_TOKEN", token)
	}, 0)

	Describe("HandleRequest", func() {
		It("invokes QuietWhiteboard.HandleInput with payload text with the right arguments", func() {
			params["text"] = "makeCoffee two sugars no milk"
			request := makeRequest(params)

			expectedContext := SlackContext{User: slackUser, Channel: slackChannel}

			handlerFunc.ServeHTTP(writer, request)

			Expect(mockWhiteBoard.HandleInputCalled).To(BeTrue())
			Expect(mockWhiteBoard.HandleInputArgs.Text).To(Equal("makeCoffee two sugars no milk"))
			Expect(mockWhiteBoard.HandleInputArgs.Context).To(Equal(expectedContext))
		})

		It("returns the JSON representation of the QuietWhiteboard response", func() {
			request := makeRequest(params)

			handlerFunc.ServeHTTP(writer, request)

			Expect(writer.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(writer.Body.String()).To(Equal(`{"text":"This is a mock message"}`))
		})

		AssertDoesNotInvokeHandleInput := func() func() {
			return func() {
				request := makeRequest(params)
				handlerFunc.ServeHTTP(writer, request)

				Expect(mockWhiteBoard.HandleInputCalled).To(BeFalse())
			}
		}

		AssertReturns403Forbidden := func() func() {
			return func() {
				request := makeRequest(params)
				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Code).To(Equal(http.StatusForbidden))
			}
		}

		AssertReturnsErrorMessage := func() func() {
			return func() {
				request := makeRequest(params)
				handlerFunc.ServeHTTP(writer, request)

				Expect(writer.Body.String()).To(Equal("Uh-oh, something went wrong... sorry!"))
			}
		}

		Context("when the SLACK_TOKEN environment variable is blank", func() {
			BeforeEach(func() {
				params["token"] = ""
				os.Unsetenv("SLACK_TOKEN")
			})

			It("does not invoke QuietWhiteboard.HandleInput", AssertDoesNotInvokeHandleInput())
			It("returns a 403 Forbidden", AssertReturns403Forbidden())
			It("returns an error message'", AssertReturnsErrorMessage())
		})

		Context("when the token is invalid", func() {
			BeforeEach(func() {
				params["token"] = "invalid"
			})

			It("does not invoke QuietWhiteboard.HandleInput", AssertDoesNotInvokeHandleInput())
			It("returns a 403 Forbidden", AssertReturns403Forbidden())
			It("returns an error message'", AssertReturnsErrorMessage())
		})
	})
})
