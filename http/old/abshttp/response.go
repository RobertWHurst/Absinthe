package abshttp

import (
	"net/http"
)

func NewResponse(request *Request) *Response {
	return &Response{
		Headers:  make(http.Header),
		Trailers: make(http.Header),
		Request:  request,

		bodyLength: 0,
		hasEnded:   false,
	}
}

type Response struct {
	Headers    http.Header
	Trailers   http.Header
	Request    *Request
	StatusCode int
	Subjects   Subjects

	bodyLength int64
	hasEnded   bool
}
