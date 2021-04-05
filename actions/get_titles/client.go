package get_titles

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type HTTPClient interface {
	SendRequest(ctx context.Context, wg *sync.WaitGroup, addr string, out chan Info)
}

type Client struct {
	transport http.RoundTripper
}

func (C *Client) SendRequest(ctx context.Context, wg *sync.WaitGroup, addr string, out chan Info) {
	defer wg.Done()
	// Setup http transport.
	client := http.Client{
		Transport: C.transport,
	}

	r, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		out <- Info{
			Url: addr,
		}
	}
	r = r.WithContext(ctx)
	response, err := client.Do(r)
	if err != nil {
		out <- Info{
			Url: addr,
		}
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		out <- Info{
			Url: addr,
			Err: fmt.Sprintf("error reading response body %v", err),
		}
	}
	if response.StatusCode < 200 || response.StatusCode > 299 || response.StatusCode == http.StatusNoContent {
		if response.StatusCode == http.StatusServiceUnavailable {
			out <- Info{
				Url: addr,
				Err: fmt.Sprintf("service unavailable"),
			}
		}
		if response.StatusCode == 500 && strings.HasPrefix(string(responseBody), "Error when calling internal api for GetApplicationStatus ") {
			out <- Info{
				Url: addr,
				Err: fmt.Sprintf("service is slow, try again"),
			}
		}
		out <- Info{
			Url: addr,
			Err: fmt.Sprintf("unexpected status code %v", response.StatusCode),
		}
	}
	title := ExtractTitle(responseBody)
	out <- Info{
		Url:   addr,
		Title: title,
	}
}

func ExtractTitle(responseBody []byte) string {
	var re = regexp.MustCompile(`(?m)<\s*title[^>]*>((.|\n)*?)<\s*\/\s*title>`)
	result := re.FindSubmatch(responseBody)
	return string(result[1])
}

func NewClient(transport http.RoundTripper) HTTPClient {
	return &Client{
		transport: transport,
	}
}
