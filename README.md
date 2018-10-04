# go-clients

Net client wrapper with convenience methods. 
Support clients `HttpClient` and `JsonClient`.

#### Clients
Create a new client.
```go
c := clients.DefaultJson()
```

#### Transformations
Transform response data to models.
```go
type Model struct {}
```
```go
func InfiniteRequests() {
    results := []*Model
    
    res, err := c.Get("http://server/api/model/id").
        Response().
        TransformTo(models, reflect.TypeOf(&Model{})).
        End()
}
```

#### Pagination
Provide a pagination function to handle response across multiple pages.
```go
func InfinitePagination(res *http.Response) (string, bool) {
    return res.Request.URL.String(), true
}
```
```go
func PaginatedRequests() {
    results := []*Model{}
    
    res, err := c.Get("http://server/api/model").
        Response().
        Paginated(InfinitePagination).
        TransformTo(&models, reflect.TypeOf(&Model{})).
        End()
}
```

#### Validation
Provide a validation function which takes a response and returns an error.
```go
func CheckStatusCodes(res *http.Response) error {
    switch res.StatusCode {
    case 200,201,202,204:
        return nil
    }
    return fmt.Errorf("Invalid status [%d].", res.StatusCode)
}
```
```go
func ValidatedRequests() error {
    model := &Model{}
    
    err := c.Post("http://server/api/model").
    	Send()
    	Response(model).
    	Validate(CheckStatusCodes)
}
```

#### Error Handling
Errors are persisted during the request / response builder. Get either the error or the response and error.
```go
res, err := c.Get("http://server/api/model").
    Reponse().
    End()
```
```go
err := c.Get("http://server/api/model").
    Reponse().
    Err()
```
