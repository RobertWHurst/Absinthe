package abshttp

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func NewRequest(request *http.Request) *Request {
	requestID := fmt.Sprint(time.Now().UnixNano(), rand.Int())
	return &Request{
		Method:    request.Method,
		URL:       *request.URL,
		Headers:   request.Header,
		Trailers:  request.Trailer,
		Subjects:  NewSubjects(request.Method, request.URL.Path, requestID),
		RequestID: requestID,
		// TODO: Figure out if there is a request body
		HasBody: false,
	}
}

type Request struct {
	Method    string
	URL       url.URL
	Headers   http.Header
	Trailers  http.Header
	Subjects  Subjects
	RequestID string
	HasBody   bool
}
