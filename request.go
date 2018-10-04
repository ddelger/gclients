package gclients

import (
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type RequestInterface interface {
	Send(interface{}) RequestInterface

	WithQuery(interface{}) RequestInterface
	WithHeader(string, string) RequestInterface

	Response() ResponseInterface
}

type Request struct {
	client *http.Client

	err     error
	url     string
	method  string
	body    io.Reader
	query   interface{}
	headers map[string]string
}

func (r *Request) Do() (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	u, err := url.Parse(r.url)
	if err != nil {
		return nil, err
	}

	if r.query != nil {
		q, err := query.Values(r.query)
		if err != nil {
			return nil, err
		}

		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(r.method, u.String(), r.body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	return r.client.Do(req)
}
