package client

import (
	"net/http"
	"net/url"
	"time"

	"h12.io/socks"
)

const timeout = time.Duration(5 * time.Second)

func New(timeout time.Duration) *http.Client {
	client := &http.Client{
		Timeout: timeout,
	}

	return client
}

func WithProxy(proxy string) (*http.Client, error) {
	if proxy == "" {
		return New(timeout), nil
	}

	url, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{Dial: socks.Dial(url.String())}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	return client, nil
}
