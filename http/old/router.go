package abs

import (
	"regexp"

	"github.com/RobertWHurst/absinthe/abshttp"
)

type Router struct {
	basePath        string
	root            bool
	conn            *Conn
	rpcHandlerFuncs map[string][]*RPCHandlerFunc
	ginHandlerFuncs map[string][]*abshttp.Handler
}

func (r *Router) Bind(rpcHandlerName string, handlers ...RPCHandlerFunc) IRoutes {
	return r.returnObj()
}

func (r *Router) Use(middleware ...abshttp.Handler) IRoutes {
	return r.returnObj()
}

func (r *Router) Any(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("*", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) DELETE(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("DELETE", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) GET(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("GET", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) HEAD(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("HEAD", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) OPTIONS(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("OPTIONS", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) PATCH(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("PATCH", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) POST(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("POST", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) PUT(relativePath string, handlers ...abshttp.Handler) IRoutes {
	r.handle("PUT", relativePath, handlers)
	return r.returnObj()
}

func (r *Router) Handle(method string, relativePath string, handlers ...abshttp.Handler) IRoutes {
	if matches, err := regexp.MatchString("^[A-Z]+$", method); !matches || err != nil {
		panic("http method " + method + " is not valid")
	}
	return r.handle(method, relativePath, handlers)
}

type IRoutes interface {
	Bind(rpcHandlerName string, handlers ...RPCHandlerFunc) IRoutes
	Use(...abshttp.Handler) IRoutes
	Handle(string, string, ...abshttp.Handler) IRoutes
	Any(string, ...abshttp.Handler) IRoutes
	GET(string, ...abshttp.Handler) IRoutes
	POST(string, ...abshttp.Handler) IRoutes
	DELETE(string, ...abshttp.Handler) IRoutes
	PATCH(string, ...abshttp.Handler) IRoutes
	PUT(string, ...abshttp.Handler) IRoutes
	OPTIONS(string, ...abshttp.Handler) IRoutes
	HEAD(string, ...abshttp.Handler) IRoutes
}

type RPCHandlerFunc func(args *interface{}, reply *interface{})

func (r *Router) handle(method, relativePath string, handlers abshttp.HandlerChain) IRoutes {
	r.conn.httpSubscribe(method, relativePath, handlers)
	return r.returnObj()
}

func (r *Router) returnObj() IRoutes {
	if r.root {
		return r.conn
	}
	return r
}
