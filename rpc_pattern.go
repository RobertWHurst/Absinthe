package absinthe

import (
	"fmt"
	"regexp"
	"strings"
)

var rpcPatternKeyPattern = regexp.MustCompile(`^\$([\w\d]+)(?:\((.*)\))?$`)
var rpcPatternWildCardPattern = regexp.MustCompile(`^\$(?:\((.*)\))?$`)
var rpcPatternChunkEscapePattern = regexp.MustCompile(`([\(\)\[\]\{\}\.\?\+\*\^\$\|\-\\])`)

type RPCPattern struct {
	PatternSrc string
	Pattern    *regexp.Regexp
}

func NewRPCPattern(patternSrc string) (*RPCPattern, error) {
	patternSrcChunks := strings.Split(patternSrc, ".")

	regExpSrcChunks := make([]string, 0)
	for _, chunk := range patternSrcChunks {
		if len(chunk) == 0 && len(patternSrcChunks) != 1 {
			return nil, fmt.Errorf("empty pattern chunk in pattern %s", patternSrc)
		}

		var regExpChunkSrc string
		if matches := rpcPatternKeyPattern.FindStringSubmatch(chunk); len(matches) == 3 {
			key := matches[1]
			subPattern := matches[2]
			if len(subPattern) == 0 {
				subPattern = `[^\.]+`
			}
			regExpChunkSrc = "(?P<" + key + ">" + subPattern + ")"

		} else if matches := rpcPatternWildCardPattern.FindStringSubmatch(chunk); len(matches) == 2 {
			subPattern := matches[1]
			if len(subPattern) == 0 {
				subPattern = `[^\.]+`
			}
			regExpChunkSrc = subPattern

		} else {
			regExpChunkSrc = rpcPatternChunkEscapePattern.ReplaceAllString(chunk, `\$1`)
		}

		regExpSrcChunks = append(regExpSrcChunks, regExpChunkSrc)
	}
	regExpSrc := `^` + strings.Join(regExpSrcChunks, `\.`) + `$`

	return &RPCPattern{
		PatternSrc: patternSrc,
		Pattern:    regexp.MustCompile(regExpSrc),
	}, nil
}

func (p *RPCPattern) Match(path string) bool {
	return p.Pattern.MatchString(path)
}

func (p *RPCPattern) FindParams(path string) (map[string]string, bool) {
	if !p.Match(path) {
		return nil, false
	}

	params := make(map[string]string)
	subExpNames := p.Pattern.SubexpNames()
	subMatches := p.Pattern.FindStringSubmatch(path)

	for i, key := range subExpNames {
		if len(key) == 0 {
			continue
		}
		params[key] = subMatches[i]
	}

	return params, true
}

func (p *RPCPattern) GobDecode(data []byte) error {
	pattern, err := NewRPCPattern(string(data))
	if err != nil {
		return err
	}
	p.Pattern = pattern.Pattern
	p.PatternSrc = pattern.PatternSrc
	return nil
}

func (p *RPCPattern) GobEncode() ([]byte, error) {
	return []byte(p.PatternSrc), nil
}

func (p *RPCPattern) String() string {
	return fmt.Sprintf("RPC(%s)", p.Pattern.String())
}
