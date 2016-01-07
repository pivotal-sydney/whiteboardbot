package main

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/model"
	"net/http"
	"os"
	"os/signal"
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

	go startHttpServer()

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

func startHttpServer() {
	http.HandleFunc("/", HealthCheckServer)
	if err := http.ListenAndServe(":" + getHealthCheckPort(), nil); err != nil {
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

func HealthCheckServer(responseWriter http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(responseWriter, "I'm alive")
}
