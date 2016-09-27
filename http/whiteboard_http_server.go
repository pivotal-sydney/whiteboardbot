package http

import (
	"encoding/json"
	"fmt"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	"io"
	"net/http"
	"os"
)

const (
	DEFAULT_PORT = "9000"
)

type WhiteboardHttpServer struct {
	SlackClient SlackClient
	Whiteboard  QuietWhiteboard
}

func (server WhiteboardHttpServer) Run() {
	server.startHttpServer()
}

func (server WhiteboardHttpServer) startHttpServer() {
	http.HandleFunc("/", server.NewHandleRequest())

	if err := http.ListenAndServe(":"+server.getHealthCheckPort(), nil); err != nil {
		fmt.Printf("ListenAndServe: %v\n", err)
	}
}

func (server WhiteboardHttpServer) getHealthCheckPort() (port string) {
	if port = os.Getenv("PORT"); len(port) == 0 {
		fmt.Printf("Warning, PORT not set. Defaulting to %+v\n", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	return
}

func (server WhiteboardHttpServer) NewHandleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		cmdArgs := req.FormValue("text")
		token := req.FormValue("token")

		if "" == os.Getenv("SLACK_TOKEN") || token != os.Getenv("SLACK_TOKEN") {
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, "Uh-oh, something went wrong... sorry!")
			return
		}

		context := server.extractSlackContext(req)

		result := server.Whiteboard.ProcessCommand(cmdArgs, context)

		resultJson, err := jsonify(result.Entry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resultJson)
	}
}

func jsonify(s fmt.Stringer) ([]byte, error) {
	response := struct {
		Text string `json:"text"`
	}{s.String()}

	return json.Marshal(response)
}

func (server WhiteboardHttpServer) extractSlackContext(req *http.Request) SlackContext {
	userId := req.FormValue("user_id")
	slackUser := server.SlackClient.GetUserDetails(userId)

	channelName := req.FormValue("channel_name")
	channelId := req.FormValue("channel_id")
	slackChannel := SlackChannel{ChannelId: channelId, ChannelName: channelName}

	return SlackContext{User: slackUser, Channel: slackChannel}
}
