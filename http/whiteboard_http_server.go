package http

import (
	"encoding/json"
	"fmt"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	"net/http"
	"os"
)

const (
	DEFAULT_PORT = "9000"
)

type WhiteboardHttpServer struct {
	Store Store
}

func NewWhiteboardHttpServer(store Store) (server WhiteboardHttpServer) {
	server = WhiteboardHttpServer{Store: store}
	return
}

func (server WhiteboardHttpServer) Run() {
	whiteboard := NewQuietWhiteboard(&RealRestClient{}, server.Store)
	server.startHttpServer(whiteboard)
}

func (server WhiteboardHttpServer) startHttpServer(wb QuietWhiteboardApp) {
	http.HandleFunc("/", server.NewHandleRequest(wb))

	if err := http.ListenAndServe(":"+server.getHealthCheckPort(), nil); err != nil {
		fmt.Printf("ListenAndServe: %v\n", err)
	}
}

func (server WhiteboardHttpServer) getHealthCheckPort() (port string) {
	if port = os.Getenv("PORT"); len(port) == 0 {
		fmt.Printf("Warning, PORT not set. Defaulting to %+v\n", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	return
}

func (server WhiteboardHttpServer) NewHandleRequest(wb QuietWhiteboard) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		cmdArgs := req.FormValue("text")
		token := req.FormValue("token")

		if "" == os.Getenv("SLACK_TOKEN") || token != os.Getenv("SLACK_TOKEN") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		response := wb.HandleInput(cmdArgs)
		j, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}
