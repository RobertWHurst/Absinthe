package absinthe

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/coreos/go-semver/semver"
)

type clientDescriptor struct {
	ID                 string
	Version            *semver.Version
	Name               string
	HandlerDescriptors []handlerDescriptor
	RouteDescriptors   []routeDescriptor
}

func newClientDescriptor(name string, version *semver.Version) clientDescriptor {
	now := time.Now()
	id := fmt.Sprintf(
		"%04d%02d%02d%02d%02d%02d%04d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond()/1000000,
	)
	for n := 0; n < 4; n++ {
		id += strconv.Itoa(n)
	}

	return clientDescriptor{
		ID:                 hex.EncodeToString([]byte(id)),
		Version:            version,
		Name:               name,
		HandlerDescriptors: make([]handlerDescriptor, 0),
		RouteDescriptors:   make([]routeDescriptor, 0),
	}
}

func (c *clientDescriptor) registerHandler(descriptor handlerDescriptor) {
	c.HandlerDescriptors = append(c.HandlerDescriptors, descriptor)
}

func (c *clientDescriptor) registerRoute(descriptor routeDescriptor) {
	c.RouteDescriptors = append(c.RouteDescriptors, descriptor)
}

func (c *clientDescriptor) testCall(namespace string, inType, outType reflect.Type) bool {
	for _, descriptor := range c.HandlerDescriptors {
		if descriptor.testCall(namespace, inType, outType) {
			return true
		}
	}
	return false
}

func (c *clientDescriptor) testRequest() bool {
	for _, descriptor := range c.RouteDescriptors {
		if descriptor.testRequest() {
			return true
		}
	}
	return false
}

type handlerDescriptor struct {
	Namespace string
	InType    reflect.Type
	OutType   reflect.Type
}

func (h *handlerDescriptor) testCall(namespace string, inType, outType reflect.Type) bool {
	return h.Namespace == namespace &&
		h.InType == inType &&
		h.OutType == outType
}

type routeDescriptor struct {
	Pattern string
	Method  string
}

func (r *routeDescriptor) testRequest() bool {
	panic("Not implemented")
}
