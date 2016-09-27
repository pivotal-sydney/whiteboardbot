package main

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/http"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	"os"
	"os/signal"
	"syscall"
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

func cleanup() {
	if redisConnectionPool != nil {
		fmt.Println("Closing Redis connection pool")
		redisConnectionPool.Close()
	}
}

func main() {
	store := RealStore{Pool: redisConnectionPool}
	slackClient := makeSlackClient()

	gateway := WhiteboardGateway{RestClient: &RealRestClient{}}
	whiteboard := NewQuietWhiteboard(gateway, &store, &RealClock{})
	server := WhiteboardHttpServer{SlackClient: slackClient, Whiteboard: whiteboard}
	server.Run()
}

func makeSlackClient() SlackClient {
	api := slack.New(os.Getenv("WB_BOT_API_TOKEN"))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	return &Slack{SlackRtm: rtm}
}
