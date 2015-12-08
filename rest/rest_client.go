package rest
import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"os"
	"errors"
)

type RestClient interface {
	Post(request WhiteboardRequest) (itemId string, ok bool)
}

type RealRestClient struct{}

func (RealRestClient) Post(request WhiteboardRequest) (itemId string, ok bool) {
	json, _ := json.Marshal(request)
	http.DefaultClient.CheckRedirect = noRedirect
	resp, err := http.Post(os.Getenv("WB_HOST_URL") + "/standups/1/items", "application/json", bytes.NewReader(json))
	fmt.Printf("Response: %v, Err: %v, json: %v", resp, err, string(json))
	ok = resp.StatusCode == http.StatusFound
	if ok {
		itemId = resp.Header.Get("Item-Id")
	}
	return
}
func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("Don't redirect!")
}