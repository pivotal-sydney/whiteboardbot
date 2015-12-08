package main
import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/entry"
	"github.com/xtreme-andleung/whiteboardbot/rest"
)

func main() {
	api := slack.New("xoxb-15808945314-Pztfx4s7YG00QAO6DlajZZdO")

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	clock := entry.RealClock{}
	restClient := rest.RealRestClient{}

	Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				go ParseMessageEvent(rtm, restClient, clock, ev)
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
			}
		}
	}
}
