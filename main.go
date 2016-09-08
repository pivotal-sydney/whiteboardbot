package main

import (
	"encoding/json"
	"fmt"
	"github.com/benjamintanweihao/slack"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	"github.com/pivotal-sydney/whiteboardbot/model"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	DEFAULT_PORT = "9000"
)

var redisConnectionPool = NewPool()

func init() {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt)
	signal.Notify(shutdownChannel, syscall.SIGTERM)
	go func() {
		<-shutdownChannel
		cleanup()
		os.Exit(1)
	}()
}

func main() {
	api := slack.New(os.Getenv("WB_BOT_API_TOKEN"))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	store := RealStore{redisConnectionPool}
	slackClient := Slack{SlackRtm: rtm}
	whiteboard := NewWhiteboard(&slackClient, &RealRestClient{}, model.RealClock{}, &store)

	go startHttpServer(whiteboard)

	// Might not even need this.
Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				go whiteboard.ParseMessageEvent(ev)
			case *slack.InvalidAuthEvent:
				fmt.Println("Invalid credentials")
				break Loop
			default:
			}
		}
	}
}

func cleanup() {
	if redisConnectionPool != nil {
		fmt.Println("Closing Redis connection pool")
		redisConnectionPool.Close()
	}
}

func startHttpServer(wb WhiteboardApp) {
	http.HandleFunc("/", NewHandleRequest(wb))
	if err := http.ListenAndServe(":"+getHealthCheckPort(), nil); err != nil {
		fmt.Printf("ListenAndServe: %v\n", err)
	}
}

func getHealthCheckPort() (port string) {
	if port = os.Getenv("PORT"); len(port) == 0 {
		fmt.Printf("Warning, PORT not set. Defaulting to %+v\n", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	return
}

type Response struct {
	Text string `json:"text"`
}

func NewHandleRequest(wb WhiteboardApp) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// fmt.Printf("--> %s\n\n", formatRequest(req))

		channelId := req.FormValue("channel_id")
		// channelName := req.FormValue("channel_name")
		// teamDomain := req.FormValue("pivotal")
		// teamId := req.FormValue("team_id")
		cmdArgs := req.FormValue("text")
		// token := req.FormValue("token")
		// userId := req.FormValue("user_id")
		// userName := req.FormValue("user_name")

		ev := slack.MessageEvent{}
		ev.Channel = channelId

		// TODO: Here we get the response
		// TODO: Don't think we would need ev.
		wb.HandleInput(cmdArgs, &ev)

		response := Response{USAGE}

		j, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}

// TODO: Remove when done. We don't need this.
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
