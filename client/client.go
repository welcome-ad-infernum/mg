package client

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"h12.io/socks"
)

const timeout = time.Duration(5 * time.Second)

func New() http.Client {
	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	return client
}

func WithProxy(proxy string) (http.Client, error) {
	if proxy == "" {
		return New(), nil
	}

	url, err := url.Parse(proxy)
	if err != nil {
		return http.Client{}, err
	}

	q := url.Query()
	q.Add("timeout", timeout.String())

	url.RawQuery = q.Encode()

	transport := &http.Transport{Dial: socks.Dial(url.String())}

	client := http.Client{
		Transport: transport,
	}

	return client, nil
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}
