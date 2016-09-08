package app

type QuietWhiteboardApp struct {
	RestClient RestClient
	CommandMap map[string]func(input string) Response
}

type Response struct {
	Text string `json:"text"`
}

// TODO: Send a message to a channel
func NewQuietWhiteboard(restClient RestClient) (whiteboard QuietWhiteboardApp) {
	whiteboard = QuietWhiteboardApp{RestClient: restClient}
	whiteboard.CommandMap = make(map[string]func(input string) Response)
	whiteboard.init()
	return
}

func (whiteboard QuietWhiteboardApp) init() {
	whiteboard.registerCommand("?", whiteboard.handleUsageCommand)
}

func (whiteboard QuietWhiteboardApp) registerCommand(
	command string,
	callback func(input string) Response) {
	whiteboard.CommandMap[command] = callback
}

func (whiteboard QuietWhiteboardApp) handleUsageCommand(_ string) Response {
	return Response{USAGE}
}

func (whiteboard QuietWhiteboardApp) HandleInput(input string) Response {
	command, input := readNextCommand(input)
	return whiteboard.handleCommand(command, input)
}

func (whiteboard QuietWhiteboardApp) handleCommand(command, input string) Response {
	for key, _ := range whiteboard.CommandMap {
		if matches(command, key) {
			callback := whiteboard.CommandMap[key]
			return callback(input)
		}
	}

	return Response{"Ooops"}
}
