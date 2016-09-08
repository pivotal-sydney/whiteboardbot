package app

type QuietWhiteboardApp struct {
	RestClient RestClient
	CommandMap map[string]func(input string) string
}

func NewQuietWhiteboard(restClient RestClient) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{RestClient: restClient}
	whiteboard.init()
	return
}

func (whiteboard QuietWhiteboardApp) init() {
	whiteboard.registerCommand("?", whiteboard.handleUsageCommand)
}

func (whiteboard QuietWhiteboardApp) registerCommand(
	command string,
	callback func(input string) string) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string) string {
	//whiteboard.SlackClient.PostMessageWithMarkdown(USAGE, ev.Channel, "")
	return USAGE
}
