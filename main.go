package main
import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/model"
	"github.com/xtreme-andleung/whiteboardbot/rest"
	"github.com/xtreme-andleung/whiteboardbot/app"
	"os"
)

func main() {
	api := slack.New(os.Getenv("WB_BOT_API_TOKEN"))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

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
				fmt.Printf("Invalid credentials")
				break Loop

			default:
			}
		}
	}
}
