package gclients

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type ClientInterface interface {
	Get(string) RequestInterface
	Put(string) RequestInterface
	Post(string) RequestInterface
	Delete(string) RequestInterface
}

type AbstractClient struct {
	client  *http.Client
	headers map[string]string
}

func Default() *AbstractClient {
	jar, _ := cookiejar.New(nil)

	return (&AbstractClient{client: &http.Client{}, headers: make(map[string]string)}).
		WithTimeout(60 * time.Second).
		WithCookieJar(jar).
		WithTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}).
		WithCheckRedirects(func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse })
}

func (c *AbstractClient) WithHeader(key, value string) *AbstractClient {
	c.headers[key] = value
	return c
}

func (c *AbstractClient) WithTimeout(timeout time.Duration) *AbstractClient {
	c.client.Timeout = timeout
	return c
}

func (c *AbstractClient) WithCookieJar(jar *cookiejar.Jar) *AbstractClient {
	c.client.Jar = jar
	return c
}

func (c *AbstractClient) WithTransport(transport http.RoundTripper) *AbstractClient {
	c.client.Transport = transport
	return c
}

func (c *AbstractClient) WithCheckRedirects(f func(req *http.Request, via []*http.Request) error) *AbstractClient {
	c.client.CheckRedirect = f
	return c
}
