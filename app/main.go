package main
import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/interfaces"
)

func main() {

	//	test(slack.New("xoxp-15805159827-15805160003-15810397457-f5d4c5b90f"))
	//
	//	test(new(interfaces.MockSlackClient))
	//	os.Args
	//
	////	os.Args
	////
	////	flag.String("type", string, "Type of whiteboard section")
	////	var input = "-type=hello world"
	//
	//	var args = "something with lots of spaces and \" quotes \" right?"
	//	split := strings.Split(args, " ")
	//	fmt.Printf("Number of strings: %d", len(split))


	api := slack.New("xoxb-15808945314-Pztfx4s7YG00QAO6DlajZZdO")

	rtm := api.NewRTM()
	go rtm.ManageConnection()


	Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				interfaces.ParseMessageEvent(rtm, ev)
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}

}
