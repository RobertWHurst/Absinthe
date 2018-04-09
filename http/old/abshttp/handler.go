package abshttp

type HandlerFunc func(c *Context)

type HandlerChain []Handler

func (h HandlerChain) Handle(c *Context) {
	// FIX ME: walk through each handler
	// currentHandler := 0

	// next := func() {

	// }

	h[0](c)
}
