package get_titles

import (
	"context"
	"crawler/domain"
	"encoding/json"
	"net/http"
	"sync"
)

type GetTitles struct {
	client HTTPClient
}

func (g GetTitles) Do(request domain.Request) (*domain.Response, error) {
	var actionRequest Request
	err := json.Unmarshal(request.ActionFields, &actionRequest)
	if err != nil {
		return nil, err
	}
	actionResponse := g.SendRequests(actionRequest)
	return &domain.Response{
		Result: actionResponse,
	}, nil
}

func (g GetTitles) SendRequests(serviceRequest Request) Response {
	var response Response
	ctx := context.Background()
	out := make(chan Info, len(serviceRequest.Urls))
	wg := &sync.WaitGroup{}
	for _, url := range serviceRequest.Urls {
		wg.Add(1)
		go g.client.SendRequest(ctx, wg, url, out)
	}
	wg.Wait()
	close(out)
	for data := range out {
		response.Info = append(response.Info, data)
	}
	return response
}

func NewService(transport http.RoundTripper) GetTitles {
	return GetTitles{
		NewClient(transport),
	}
}
