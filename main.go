package main
import (
	"fmt"
	"github.com/nlopes/slack"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/xtreme-andleung/whiteboardbot/entry"
)

func main() {
	api := slack.New("xoxb-15808945314-Pztfx4s7YG00QAO6DlajZZdO")

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	clock := entry.RealClock{}

	Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				ParseMessageEvent(rtm, clock, ev)
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
			}
		}
	}

}
