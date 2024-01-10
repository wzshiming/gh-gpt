package api

import (
	"fmt"
	"net/http"

	"github.com/wzshiming/gh-gpt/pkg/cache"
)

type Option func(*Client)

func WithTokenCache(c cache.Cache) func(*Client) {
	return func(client *Client) {
		client.tokenCache = c
	}
}

func WithHTTPClient(c *http.Client) func(*Client) {
	return func(client *Client) {
		client.client = c
	}
}

type Client struct {
	client     *http.Client
	tokenCache cache.Cache
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		client:     &http.Client{},
		tokenCache: cache.NewMemoryCache(),
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

type statusError struct {
	StatusCode   int
	Status       string
	ErrorMessage string
}

func (e statusError) Error() string {
	switch {
	case e.Status != "" && e.ErrorMessage != "":
		return fmt.Sprintf("%s: %s", e.Status, e.ErrorMessage)
	case e.Status != "":
		return e.Status
	case e.ErrorMessage != "":
		return e.ErrorMessage
	default:
		return "something went wrong"
	}
}
