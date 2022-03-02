package client

import (
	"net"
	"net/http"
	"time"
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

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}
