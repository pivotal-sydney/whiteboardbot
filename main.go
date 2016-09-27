package main

import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	. "github.com/pivotal-sydney/whiteboardbot/http"
	. "github.com/pivotal-sydney/whiteboardbot/model"
	. "github.com/pivotal-sydney/whiteboardbot/slack"
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

	rtm := makeSlackRTM()
	slackClient := Slack{SlackRtm: rtm}

	gateway := WhiteboardGateway{RestClient: &RealRestClient{}}
	whiteboard := NewQuietWhiteboard(gateway, &store, &RealClock{})

	httpServer := WhiteboardHttpServer{SlackClient: &slackClient, Whiteboard: whiteboard}
	go httpServer.Run()

	slackBotServer := SlackBotServer{SlackClient: &slackClient, Whiteboard: whiteboard}
	slackBotServer.Run(rtm)
}

func makeSlackRTM() (rtm *slack.RTM) {
	api := slack.New(os.Getenv("WB_BOT_API_TOKEN"))
	rtm = api.NewRTM()
	go rtm.ManageConnection()
	return
}
