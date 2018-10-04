package gclients

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type HttpClient struct {
	*AbstractClient
}

func Http(client *AbstractClient) ClientInterface {
	return &HttpClient{AbstractClient: client}
}

func (c *HttpClient) Get(url string) RequestInterface {
	return c.Request(url, http.MethodGet)
}

func (c *HttpClient) Put(url string) RequestInterface {
	return c.Request(url, http.MethodPut)
}

func (c *HttpClient) Post(url string) RequestInterface {
	return c.Request(url, http.MethodPost)
}

func (c *HttpClient) Delete(url string) RequestInterface {
	return c.Request(url, http.MethodDelete)
}

func (c *HttpClient) Request(url, method string) *HttpRequest {
	r := &HttpRequest{Request: &Request{url: url, method: method, client: c.client, headers: make(map[string]string)}}

	for k, v := range c.headers {
		r.headers[k] = v
	}

	return r
}

type HttpRequest struct {
	*Request
}

func (r *HttpRequest) Send(v interface{}) RequestInterface {
	if b, ok := v.(io.Reader); !ok {
		r.err = fmt.Errorf("Invalid body type [%s]. Must be of type io.Reader.", reflect.TypeOf(v))
	} else {
		r.body = b
	}
	return r
}

func (r *HttpRequest) WithQuery(query interface{}) RequestInterface {
	r.query = query
	return r
}

func (r *HttpRequest) WithHeader(key, value string) RequestInterface {
	r.headers[key] = value
	return r
}

func (r *HttpRequest) Response() ResponseInterface {
	res, err := r.Do()
	return &HttpResponse{&Response{res: res, err: err}}
}

type HttpResponse struct {
	*Response
}

func (r *HttpResponse) Paginated(f PaginationFunc) ResponseInterface {
	return r
}

func (r *HttpResponse) TransformTo(v interface{}, t reflect.Type) ResponseInterface {
	return r
}
