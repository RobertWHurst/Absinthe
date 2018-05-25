package absinthe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var restMethods = []string{
	"post",
	"get",
	"put",
	"patch",
	"delete",
	"head",
	"options",
}

var restPatternAndPaths = map[string][2][]string{
	"": [2][]string{
		[]string{
			"",
			"/",
		},
		[]string{
			"/alph",
			"/alpha",
			"/alpha/",
			"/1",
			"/1/",
		},
	},
	"/alpha": [2][]string{
		[]string{
			"/alpha",
			"/alpha/",
		},
		[]string{
			"",
			"/",
			"/alph",
			"/alpha/beta",
			"/alpha/beta/",
			"/1",
			"/1/",
			"/1/2",
			"/1/2/",
		},
	},
	"/alpha/beta": [2][]string{
		[]string{
			"/alpha/beta",
			"/alpha/beta/",
		},
		[]string{
			"",
			"/",
			"/alph",
			"/alpha",
			"/alpha/",
			"/alpha/beta/gamma",
			"/alpha/beta/gamma/",
		},
	},
	"/*": [2][]string{
		[]string{
			"/alpha",
			"/beta",
			"/gama",
		},
		[]string{
			"",
			"/",
			"/alpha/beta",
			"/alpha/beta/",
			"/1/2",
			"/1/2/",
		},
	},
	"/alpha/*": [2][]string{
		[]string{
			"/alpha/alpha",
			"/alpha/beta",
			"/alpha/gamma",
			"/alpha/1",
			"/alpha/2",
			"/alpha/3",
		},
		[]string{
			"/alpha/beta/gamma",
			"/alpha/beta/gamma/",
			"/1/2",
			"/1/2/",
		},
	},
	"/alpha/*/gamma": [2][]string{
		[]string{
			"/alpha/alpha/gamma",
			"/alpha/beta/gamma",
			"/alpha/gamma/gamma",
			"/alpha/1/gamma",
			"/alpha/2/gamma",
			"/alpha/3/gamma",
		},
		[]string{
			"/alpha/beta",
			"/alpha/beta/",
			"/alpha/beta/gamma/delta",
			"/alpha/beta/gamma/delta/",
			"/alpha/beta/1",
			"/alpha/beta/1/",
		},
	},
	`/*(\d\w{3})`: [2][]string{
		[]string{
			"/1abc",
			"/2efg",
		},
		[]string{
			"/aabc",
			"/1ab",
		},
	},
	"/:key": [2][]string{
		[]string{
			"/alpha",
			"/beta",
			"/gama",
		},
		[]string{
			"",
			"/",
			"/alpha/beta",
			"/alpha/beta/",
			"/1/2",
			"/1/2/",
		},
	},
	"/alpha/:key": [2][]string{
		[]string{
			"/alpha/alpha",
			"/alpha/beta",
			"/alpha/gamma",
			"/alpha/1",
			"/alpha/2",
			"/alpha/3",
		},
		[]string{
			"/alpha/beta/gamma",
			"/alpha/beta/gamma/",
			"/1/2",
			"/1/2/",
		},
	},
	"/alpha/:key/gamma": [2][]string{
		[]string{
			"/alpha/alpha/gamma",
			"/alpha/beta/gamma",
			"/alpha/gamma/gamma",
			"/alpha/1/gamma",
			"/alpha/2/gamma",
			"/alpha/3/gamma",
		},
		[]string{
			"/alpha/beta",
			"/alpha/beta/",
			"/alpha/beta/gamma/delta",
			"/alpha/beta/gamma/delta/",
			"/alpha/beta/1",
			"/alpha/beta/1/",
		},
	},
	`/:key(\d\w{3})`: [2][]string{
		[]string{
			"/1abc",
			"/2efg",
		},
		[]string{
			"/aabc",
			"/1ab",
		},
	},
	`/?/+/(/)/./-/_/{}`: [2][]string{
		[]string{
			"/?/+/(/)/./-/_/{}",
		},
		[]string{
			"/////a/_",
			"",
			"/",
		},
	},
	`/?/+/(/)/./-/_/*(a?b+c{2})`: [2][]string{
		[]string{
			"/?/+/(/)/./-/_/abcc",
			"/?/+/(/)/./-/_/bbbbcc",
		},
		[]string{
			"/?/+/(/)/./-/_/abc",
			"/?/+/(/)/./-/_",
			"/////a/_",
			"",
			"/",
		},
	},
}

func TestNewRESTRoute(t *testing.T) {
	for _, method := range restMethods {
		for pattern := range restPatternAndPaths {
			_, err := NewRESTRoute(method, pattern)
			assert.NoErrorf(t, err, "Failed to create route for method %s and pattern %s", method, pattern)
		}
	}
}

func TestRESTRouteMatch(t *testing.T) {
	for _, method := range restMethods {
		for pattern, paths := range restPatternAndPaths {
			route, err := NewRESTRoute(method, pattern)
			assert.NoError(t, err, "Failed to create route for method %s and pattern %s", method, pattern)

			for _, validPath := range paths[0] {
				assert.True(t, route.Match(method, validPath), "route with method %s and path %s should match method %s and path %s but it does not", method, pattern, method, validPath)
			}
			for _, invalidPath := range paths[1] {
				assert.False(t, route.Match(method, invalidPath), "route with method %s and path %s should not match method %s and path %s but it does", method, pattern, method, invalidPath)
			}
		}
	}
}

func TestRESTRouteFindParams(t *testing.T) {
	route, err := NewRESTRoute("get", "/1/:alpha/2/:beta/3/:gamma")
	assert.NoError(t, err, "Failed to create route for testing FindParams")

	params, ok := route.FindParams("get", "/1/one/2/two/3/three")
	assert.True(t, ok, "FindParams did find the path or method matching")

	assert.Equal(t, 3, len(params))
	assert.Equal(t, "one", params["alpha"])
	assert.Equal(t, "two", params["beta"])
	assert.Equal(t, "three", params["gamma"])

	_, ok = route.FindParams("post", "/1/one/2/two/3/three")
	assert.False(t, ok, "A get route should not match a post request")
}
