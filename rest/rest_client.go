package rest
import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"os"
	"errors"
	"strings"
	"github.com/xtreme-andleung/whiteboardbot/model"
)

type RestClient interface {
	Post(request model.WhiteboardRequest) (itemId string, ok bool)
}

type RealRestClient struct{}

func (RealRestClient) Post(request model.WhiteboardRequest) (itemId string, ok bool) {
	json, _ := json.Marshal(request)
	fmt.Printf("Posting entry to whiteboard:\n%v\n", string(json))
	http.DefaultClient.CheckRedirect = noRedirect
	url := os.Getenv("WB_HOST_URL")
	standupId := os.Getenv("WB_STANDUP_IP")
	if len(request.Id) > 0 {
		url += "/items/" + request.Id
	} else {
		url += "/standups/" + standupId + "/items"
	}
	httpRequest, err := http.NewRequest(toHttpVerb(request.Method), url, bytes.NewReader(json))
	httpRequest.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpRequest)
	fmt.Printf("Whitebord Response: %v, Err: %v\n, Url: %v\n", resp, err, url)

	ok = resp !=nil && resp.StatusCode == http.StatusFound
	if ok {
		itemId = resp.Header.Get("Item-Id")
		if (len(itemId) == 0) {
			ok = false
		}
	} else {
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