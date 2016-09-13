package main

import (
	"encoding/json"
	"fmt"
	// "github.com/nlopes/slack"
	. "github.com/pivotal-sydney/whiteboardbot/app"
	// "github.com/pivotal-sydney/whiteboardbot/model"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	DEFAULT_PORT = "9000"
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

func main() {
	// TODO: need the store
	// store := RealStore{redisConnectionPool}
	// whiteboard := NewWhiteboard(&slackClient, &RealRestClient{}, model.RealClock{}, &store)
	whiteboard := NewQuietWhiteboard(&RealRestClient{})
	startHttpServer(whiteboard)
}

func cleanup() {
	if redisConnectionPool != nil {
		fmt.Println("Closing Redis connection pool")
		redisConnectionPool.Close()
	}
}

func startHttpServer(wb QuietWhiteboardApp) {
	http.HandleFunc("/", NewHandleRequest(wb))

	if err := http.ListenAndServe(":"+getHealthCheckPort(), nil); err != nil {
		fmt.Printf("ListenAndServe: %v\n", err)
	}
}

func getHealthCheckPort() (port string) {
	if port = os.Getenv("PORT"); len(port) == 0 {
		fmt.Printf("Warning, PORT not set. Defaulting to %+v\n", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	return
}

func NewHandleRequest(wb QuietWhiteboardApp) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		cmdArgs := req.FormValue("text")
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
