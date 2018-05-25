package absinthe

type RESTContext struct {
	URL    string
	Method string
	Params map[string]string
	Next   func()
}

func (c *RESTContext) Status(statusCode int) *RESTContext {
	return c
}

func (c *RESTContext) End() {
}
