package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pivotal-sydney/whiteboardbot/model"
	"net/http"
	"os"
	"strings"
)

type RestClient interface {
	Post(request model.WhiteboardRequest) (itemId string, ok bool)
}

type RealRestClient struct{}

func (RealRestClient) Post(request model.WhiteboardRequest) (itemId string, ok bool) {
	json, _ := json.Marshal(request)
	http.DefaultClient.CheckRedirect = noRedirect
	url := os.Getenv("WB_HOST_URL")
	if len(request.Id) > 0 {
		url += "/items/" + request.Id
	} else {
		url += "/standups/94/items"
	}
	httpRequest, err := http.NewRequest(toHttpVerb(request.Method), url, bytes.NewReader(json))
	httpRequest.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpRequest)
	fmt.Printf("\nResponse: %v, Err: %v, json: %v", resp, err, string(json))
	fmt.Printf("\nURL %v", url)

	ok = resp != nil && resp.StatusCode == http.StatusFound
	if ok {
		itemId = resp.Header.Get("Item-Id")
	}
	if len(itemId) == 0 {
		itemId = request.Id
	}
	return
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("Don't redirect!")
}

func toHttpVerb(method string) (httpVerb string) {
	if len(method) > 0 {
		httpVerb = strings.ToUpper(method)
	} else {
		httpVerb = "POST"
	}
	return
}