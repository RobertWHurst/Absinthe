package absinthe

type RESTRouter struct {
	client    *Client
	baseRoute RESTRoute
	layers    []RESTRouterLayer
}

func NewRESTRouter() *RESTRouter {
	return &RESTRouter{
		layers: make([]RESTRouterLayer, 0),
	}
}

func (r *RESTRouter) Mount(path string, router *RESTRouter) error {
	if len(path) == 0 || path[len(path)-1:] != "+" {
		path += "+"
	}
	route, err := NewRESTRoute("all", path)
	if err != nil {
		return err
	}
	router.client = r.client
	r.layers = append(r.layers, RESTRouterLayer{
		Route:  route,
		Router: router,
	})
	return nil
}

func (r *RESTRouter) All(path string, middleware RESTHandler) error {
	route, err := NewRESTRoute("all", path)
	if err != nil {
		return err
	}
	r.layers = append(r.layers, RESTRouterLayer{
		Route:   route,
		Handler: middleware,
	})
	r.lazyPublishRoutes()
	return nil
}

func (r *RESTRouter) Exec(context *RESTContext) {
	currentLayerIndex := 0
	context.Next = func() {
		if currentLayerIndex < len(r.layers) {
			layer := r.layers[currentLayerIndex]
			currentLayerIndex++
			layer.Exec(context)
		} else {
			context.Status(404).End()
		}
	}
}

func (r *RESTRouter) lazyPublishRoutes() {
	// for _, layer := range r.layers {

	// }
}

type RESTRouterLayer struct {
	Route   *RESTRoute
	Router  *RESTRouter
	Handler RESTHandler
}

func (r *RESTRouterLayer) Exec(context *RESTContext) {
	params, ok := r.Route.FindParams(context.Method, context.URL)
	if !ok {
		context.Next()
		return
	}
	context.Params = params
	if r.Handler != nil {
		r.Handler(context)
	} else {
		r.Router.Exec(context)
	}
}
