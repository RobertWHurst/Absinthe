package absinthe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rpcPatternAndPaths = map[string][2][]string{
	"": [2][]string{
		[]string{
			"",
		},
		[]string{
			"alph",
			"alpha",
			"1",
		},
	},
	"alpha": [2][]string{
		[]string{
			"alpha",
		},
		[]string{
			"",
			"alph",
			"alpha.beta",
			"alpha.beta",
			"1",
			"1.2",
		},
	},
	"alpha.beta": [2][]string{
		[]string{
			"alpha.beta",
		},
		[]string{
			"",
			"alph",
			"alpha",
			"alpha.beta.gamma",
		},
	},
	"$": [2][]string{
		[]string{
			"alpha",
			"beta",
			"gama",
		},
		[]string{
			"",
			"alpha.beta",
			"1.2",
		},
	},
	"alpha.$": [2][]string{
		[]string{
			"alpha.alpha",
			"alpha.beta",
			"alpha.gamma",
			"alpha.1",
			"alpha.2",
			"alpha.3",
		},
		[]string{
			"alpha.beta.gamma",
			"1.2",
		},
	},
	"alpha.$.gamma": [2][]string{
		[]string{
			"alpha.alpha.gamma",
			"alpha.beta.gamma",
			"alpha.gamma.gamma",
			"alpha.1.gamma",
			"alpha.2.gamma",
			"alpha.3.gamma",
		},
		[]string{
			"alpha.beta",
			"alpha.beta.",
			"alpha.beta.gamma.delta",
			"alpha.beta.1",
		},
	},
	`$(\d\w{3})`: [2][]string{
		[]string{
			"1abc",
			"2efg",
		},
		[]string{
			"aabc",
			"1ab",
		},
	},
	"$key": [2][]string{
		[]string{
			"alpha",
			"beta",
			"gama",
		},
		[]string{
			"",
			"alpha.beta",
			"1.2",
		},
	},
	"alpha.$key": [2][]string{
		[]string{
			"alpha.alpha",
			"alpha.beta",
			"alpha.gamma",
			"alpha.1",
			"alpha.2",
			"alpha.3",
		},
		[]string{
			"alpha.beta.gamma",
			"1.2",
		},
	},
	"alpha.$key.gamma": [2][]string{
		[]string{
			"alpha.alpha.gamma",
			"alpha.beta.gamma",
			"alpha.gamma.gamma",
			"alpha.1.gamma",
			"alpha.2.gamma",
			"alpha.3.gamma",
		},
		[]string{
			"alpha.beta",
			"alpha.beta.gamma.delta",
			"alpha.beta.1",
		},
	},
	`$key(\d\w{3})`: [2][]string{
		[]string{
			"1abc",
			"2efg",
		},
		[]string{
			"aabc",
			"1ab",
		},
	},
	`?.+.(.).-._.{}`: [2][]string{
		[]string{
			"?.+.(.).-._.{}",
		},
		[]string{
			"....a._",
			"",
		},
	},
	`?.+.(.).-._.$(a?b+c{2})`: [2][]string{
		[]string{
			"?.+.(.).-._.abcc",
			"?.+.(.).-._.bbbbcc",
		},
		[]string{
			"?.+.(.).-._.abc",
			"?.+.(.).-._",
			"a._",
			"",
		},
	},
}

func TestNewRPCPattern(t *testing.T) {
	for pattern := range rpcPatternAndPaths {
		_, err := NewRPCPattern(pattern)
		assert.NoErrorf(t, err, "Failed to create pattern %s", pattern)
	}
}

func TestRPCPatternMatch(t *testing.T) {
	for pattern, paths := range rpcPatternAndPaths {
		pattern, err := NewRPCPattern(pattern)
		assert.NoError(t, err, "Failed to create pattern %s", pattern)

		for _, validPath := range paths[0] {
			assert.True(t, pattern.Match(validPath), "pattern %s should match path %s but it does not", pattern, validPath)
		}
		for _, invalidPath := range paths[1] {
			assert.False(t, pattern.Match(invalidPath), "pattern %s should not match path %s but it does", pattern, invalidPath)
		}
	}
}

func TestRPCPatternFindParams(t *testing.T) {
	pattern, err := NewRPCPattern("1.$alpha.2.$beta.3.$gamma")
	assert.NoError(t, err, "Failed to create pattern for testing FindParams")

	params, ok := pattern.FindParams("1.one.2.two.3.three")
	assert.True(t, ok, "FindParams did find the path or method matching")

	assert.Equal(t, 3, len(params))
	assert.Equal(t, "one", params["alpha"])
	assert.Equal(t, "two", params["beta"])
	assert.Equal(t, "three", params["gamma"])
}
