package rest

type RestClient interface {
	Post(request WhiteboardRequest) string
}

type RealRestClient struct {

}

func (client RealRestClient) Post(request WhiteboardRequest) {
}
