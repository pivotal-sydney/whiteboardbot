package app
import
(
	. "github.com/xtreme-andleung/whiteboardbot/model"
	"encoding/json"
	"fmt"
)

type WhiteboardApp struct {
	SlackClient SlackClient
	RestClient  RestClient
	Clock       Clock
	Store       Store
	EntryMap    map[string]EntryType
}

func (whiteboard WhiteboardApp) registerStandup(standupId string, channel string) {
	standup, ok := whiteboard.RestClient.GetStandup(standupId)
	if !ok {
		handleStandupNotFound(whiteboard.SlackClient, standupId, channel)
		return
	}
	standupJson, _ := json.Marshal(standup)
	whiteboard.Store.Set(channel, string(standupJson))
	whiteboard.SlackClient.PostMessage(fmt.Sprintf("Standup %v has been registered! You can now start creating Whiteboard entries!", standup.Title), channel, THUMBS_UP)
}

func handleMissingTitle(whiteboard WhiteboardApp, channel string) {
	whiteboard.SlackClient.PostMessageWithMarkdown("Hey, next time add a title along with your entry!\nLike this: `wb i My title`", channel, THUMBS_DOWN)
}
