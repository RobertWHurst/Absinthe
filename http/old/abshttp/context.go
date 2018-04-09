package abshttp

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/nats-io/go-nats"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Params  Params
	Keys    map[string]interface{}
	Errors  []*Error

	handlers     HandlersChain
	handlerIndex uint
	isAborted    bool
}

func ContextFromHTTP(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Request: request,
		Writer:  writer,
	}
}

func ContextFromNats(natsEncodedConnection *nats.EncodedConn) *Context {
	return &Context{
		Request: request,
		Writer:  writer,
	}
}

func (c *Context) Copy() *Context {}
func (c *Context) HandlerName() string {
	return runtime.FuncForPC(reflect.ValueOf(c.Handler()).Pointer()).Name()
}
func (c *Context) Handler() HandlerFunc {
	return c.handlers.Last()
}

func (c *Context) Next() {

}
