package absinthe

import "reflect"

type clientDescriptor struct {
	name               string
	handlerDescriptors []handlerDescriptor
	routeDescriptors   []routeDescriptor
}

type handlerDescriptor struct {
	name string
	in   reflect.Type
	out  reflect.Type
}

type routeDescriptor struct {
	pattern string
	method  string
}
