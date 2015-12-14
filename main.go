package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/model"
	"github.com/xtreme-andleung/whiteboardbot/rest"
	"os"
	"net/http"
)

const (
	DEFAULT_PORT = "9000"
)

func main() {
	api := slack.New(os.Getenv("WB_BOT_API_TOKEN"))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	go startHttpServer()

	clock := model.RealClock{}
	restClient := rest.RealRestClient{}

	Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				go app.ParseMessageEvent(rtm, restClient, clock, ev)
			case *slack.InvalidAuthEvent:
				fmt.Println("Invalid credentials")
				break Loop
			default:
			}
		}
	}
}

func startHttpServer() {
	http.HandleFunc("/", HealthCheckServer)
	if err := http.ListenAndServe(":" + getHealthCheckPort(), nil); err != nil {
		fmt.Printf("ListenAndServe: %v\n", err)
	}
}

func getHealthCheckPort() (port string){
	if port = os.Getenv("PORT"); len(port) == 0 {
		fmt.Printf("Warning, PORT not set. Defaulting to %+v\n", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	return
}

func HealthCheckServer(responseWriter http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(responseWriter, "I'm alive")
}
