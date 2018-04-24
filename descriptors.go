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
	id                 string
	version            *semver.Version
	name               string
	handlerDescriptors []handlerDescriptor
	routeDescriptors   []routeDescriptor
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
		id:      hex.EncodeToString([]byte(id)),
		version: version,
		name:    name,
	}
}

func (c *clientDescriptor) registerHandler(descriptor handlerDescriptor) {
	c.handlerDescriptors = append(c.handlerDescriptors, descriptor)
}

func (c *clientDescriptor) registerRoute(descriptor routeDescriptor) {
	c.routeDescriptors = append(c.routeDescriptors, descriptor)
}

func (c *clientDescriptor) testCall(namespace string, inType, outType reflect.Type) bool {
	for _, descriptor := range c.handlerDescriptors {
		if descriptor.testCall(namespace, inType, outType) {
			return true
		}
	}
	return false
}

func (c *clientDescriptor) testRequest() bool {
	for _, descriptor := range c.routeDescriptors {
		if descriptor.testRequest() {
			return true
		}
	}
	return false
}

type handlerDescriptor struct {
	namespace string
	inType    reflect.Type
	outType   reflect.Type
}

func (h *handlerDescriptor) testCall(namespace string, inType, outType reflect.Type) bool {
	return h.namespace == namespace &&
		h.inType == inType &&
		h.outType == outType
}

type routeDescriptor struct {
	pattern string
	method  string
}

func (r *routeDescriptor) testRequest() bool {
	panic("Not implemented")
}
