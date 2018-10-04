package gclients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

type JsonClient struct {
	*AbstractClient
}

func DefaultJson() ClientInterface {
	return Json(Default())
}

func Json(client *AbstractClient) ClientInterface {
	return &JsonClient{AbstractClient: client}
}

func (c *JsonClient) Get(url string) RequestInterface {
	return c.Request(url, http.MethodGet).WithHeader("Accept", "application/json")
}

func (c *JsonClient) Put(url string) RequestInterface {
	return c.Request(url, http.MethodPut).WithHeader("Content-Type", "application/json")
}

func (c *JsonClient) Post(url string) RequestInterface {
	return c.Request(url, http.MethodPost).WithHeader("Content-Type", "application/json")
}

func (c *JsonClient) Delete(url string) RequestInterface {
	return c.Request(url, http.MethodDelete)
}

func (c *JsonClient) Request(url, method string) *JsonRequest {
	r := &JsonRequest{Request: &Request{url: url, method: method, client: c.client, headers: make(map[string]string)}}

	for k, v := range c.headers {
		r.headers[k] = v
	}

	return r
}

type JsonRequest struct {
	*Request
}

func (r *JsonRequest) Send(v interface{}) RequestInterface {
	b, err := json.Marshal(v)
	if err != nil {
		r.err = err
	} else {
		r.body = bytes.NewReader(b)
	}
	return r
}

func (r *JsonRequest) WithQuery(query interface{}) RequestInterface {
	r.query = query
	return r
}

func (r *JsonRequest) WithHeader(key, value string) RequestInterface {
	r.headers[key] = value
	return r
}

func (r *JsonRequest) Response() ResponseInterface {
	res, err := r.Do()
	return &JsonResponse{Response: &Response{res: res, err: err}, request: r.Request}
}

type JsonResponse struct {
	*Response

	request    *Request
	pagination PaginationFunc
}

func (r *JsonResponse) Paginated(f PaginationFunc) ResponseInterface {
	r.pagination = f
	return r
}

func (r *JsonResponse) TransformTo(v interface{}, t reflect.Type) ResponseInterface {
	if r.err != nil {
		return r
	}

	if r.pagination == nil {
		r.Unmarshal(v)
	} else {
		if reflect.ValueOf(v).Elem().Kind() != reflect.Slice {
			r.err = fmt.Errorf("Invalid value kind [%s]. Must be of type slice.", reflect.ValueOf(v).Elem().Kind())
			return r
		}

		s := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
		p := reflect.New(s.Type())

		p.Elem().Set(s)
		r.Pagination(v, p.Interface())
	}

	return r
}

func (r *JsonResponse) Unmarshal(v interface{}) {
	defer r.res.Body.Close()

	if b, err := ioutil.ReadAll(r.res.Body); err != nil {
		r.err = err
	} else if err := json.Unmarshal(b, v); err != nil {
		r.err = err
	}
}

func (r *JsonResponse) Pagination(v, p interface{}) {
	r.Unmarshal(p)

	if r.err != nil {
		return
	}

	pe := reflect.ValueOf(p).Elem()
	ve := reflect.ValueOf(v).Elem()
	for i := 0; i < pe.Len(); i++ {
		ve.Set(reflect.Append(ve, pe.Index(i)))
	}

	if next, ok := r.pagination(r.res); ok {
		r.request.url = next
		r.request.query = nil

		res, err := r.request.Do()
		r.res = res
		r.err = err

		r.Pagination(v, p)
	}
}
