package absinthe

import (
	"fmt"
	"regexp"
	"strings"
)

var restRouteKeyPattern = regexp.MustCompile(`^:([\w\d]+)(?:\((.*)\))?$`)
var restRouteWildCardPattern = regexp.MustCompile(`^\*(?:\((.*)\))?$`)
var restRouteChunkEscapePattern = regexp.MustCompile(`([\(\)\[\]\{\}\.\?\+\*\^\$\|\-\\])`)

type RESTRoute struct {
	Method     string
	PatternSrc string
	Pattern    *regexp.Regexp
}

func NewRESTRoute(method, patternSrc string) (*RESTRoute, error) {
	selfTerminating := true
	if len(patternSrc) != 0 && patternSrc[len(patternSrc)-1:] == "+" {
		patternSrc = patternSrc[:len(patternSrc)-1]
		selfTerminating = false
	}

	patternSrcChunks := strings.Split(patternSrc, "/")

	regExpSrcChunks := make([]string, 0)
	for _, chunk := range patternSrcChunks {
		if len(chunk) == 0 {
			continue
		}

		var regExpChunkSrc string
		if matches := restRouteKeyPattern.FindStringSubmatch(chunk); len(matches) == 3 {
			key := matches[1]
			subPattern := matches[2]
			if len(subPattern) == 0 {
				subPattern = `[^/]+`
			}
			regExpChunkSrc = "(?P<" + key + ">" + subPattern + ")"

		} else if matches := restRouteWildCardPattern.FindStringSubmatch(chunk); len(matches) == 2 {
			subPattern := matches[1]
			if len(subPattern) == 0 {
				subPattern = `[^/]+`
			}
			regExpChunkSrc = subPattern

		} else {
			regExpChunkSrc = restRouteChunkEscapePattern.ReplaceAllString(chunk, `\$1`)
		}

		regExpSrcChunks = append(regExpSrcChunks, regExpChunkSrc)
	}
	regExpSrc := `^/?` + strings.Join(regExpSrcChunks, `/+`) + `/?`
	if selfTerminating {
		regExpSrc += `$`
	}

	fmt.Println(regExpSrc)

	return &RESTRoute{
		Method:     strings.ToLower(method),
		PatternSrc: patternSrc,
		Pattern:    regexp.MustCompile(regExpSrc),
	}, nil
}

func (r *RESTRoute) Match(method, path string) bool {
	return (r.Method == "all" || strings.ToLower(method) == r.Method) && r.Pattern.MatchString(path)
}

func (r *RESTRoute) FindParams(method, path string) (map[string]string, bool) {
	if !r.Match(method, path) {
		return nil, false
	}

	params := make(map[string]string)
	subExpNames := r.Pattern.SubexpNames()
	subMatches := r.Pattern.FindStringSubmatch(path)

	for i, key := range subExpNames {
		if len(key) == 0 {
			continue
		}
		params[key] = subMatches[i]
	}

	return params, true
}

func (r *RESTRoute) GobDecode(data []byte) error {
	methodData := strings.TrimSpace(string(data[:10]))
	patternData := strings.TrimSpace(string(data[10:]))
	route, err := NewRESTRoute(methodData, patternData)
	if err != nil {
		return err
	}
	r.Method = route.Method
	r.Pattern = route.Pattern
	r.PatternSrc = route.PatternSrc
	return nil
}

func (r *RESTRoute) GobEncode() ([]byte, error) {
	method := r.Method
	for len(method) < 10 {
		method += " "
	}
	return []byte(method + r.PatternSrc), nil
}

func (r *RESTRoute) String() string {
	return fmt.Sprintf("REST(%s:%s)", r.Method, r.Pattern.String())
}
