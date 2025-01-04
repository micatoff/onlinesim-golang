package onlinesim

import (
	"errors"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://onlinesim.io/api/"
)

var errWrongApiKey = errors.New("wrong api key")

type Client struct {
	ApiKey string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	c.httpClient.Transport = transport
	return nil
}
