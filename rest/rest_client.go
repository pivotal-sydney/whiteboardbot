package rest

type RestClient interface {
	Post(request WhiteboardRequest)
}

type RealRestClient struct {

}

func (client *RealRestClient) Post(request WhiteboardRequest) {
}
