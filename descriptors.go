package absinthe

type clientDescriptor struct {
	name               string
	handlerDescriptors []handlerDescriptor
	routeDescriptors   []routeDescriptor
}

type handlerDescriptor struct {
	name string
	in   Type
	out  Type
}

type routeDescriptor struct {
	pattern string
	method  string
}
