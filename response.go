package gclients

import (
	"net/http"
	"reflect"
)

type ValidationFunc func(*http.Response) error
type PaginationFunc func(*http.Response) (string, bool)

type ResponseInterface interface {
	End() (*http.Response, error)
	Err() error

	Validate(ValidationFunc) error
	Paginated(PaginationFunc) ResponseInterface
	TransformTo(interface{}, reflect.Type) ResponseInterface
}

type Response struct {
	err error
	res *http.Response
}

func (r *Response) End() (*http.Response, error) {
	return r.res, r.err
}

func (r *Response) Err() error {
	return r.err
}

func (r *Response) Validate(f ValidationFunc) error {
	r.err = f(r.res)
	return r.err
}
