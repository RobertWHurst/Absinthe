package absinthe

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/coreos/go-semver/semver"
)

// Peer is used to represent a peer absinthe instance. It has methods for
// dispatching requests to that peer over nats.
type Peer struct {
	ID          string
	Version     *semver.Version
	Name        string
	RPCPatterns map[string]RPCPattern
	RESTRoutes  map[string]RESTRoute
}

func NewPeer(name string, version *semver.Version, id string) (*Peer, error) {
	if len(id) == 0 {
		now := time.Now()
		id = fmt.Sprintf(
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
			id += strconv.Itoa(int(rand.Int31n(9)))
		}
	} else if len(id) != 22 {
		return nil, errors.New("id must be 22 numeric chars if provided")
	}

	return &Peer{
		ID:          hex.EncodeToString([]byte(id)),
		Version:     version,
		Name:        name,
		RPCPatterns: make(map[string]RPCPattern),
		RESTRoutes:  make(map[string]RESTRoute),
	}, nil
}

func (p *Peer) AddRPCPattern(pattern RPCPattern) {
	p.RPCPatterns[pattern.String()] = pattern
}

func (p *Peer) AddRESTRoute(route RESTRoute) {
	p.RESTRoutes[route.String()] = route
}

func (p *Peer) HasRPCHandlerFor(path string) bool {
	for _, knownPattern := range p.RPCPatterns {
		if knownPattern.Match(path) {
			return true
		}
	}
	return false
}

func (p *Peer) HasRestHandlerFor(method, path string) bool {
	for _, knownRoute := range p.RESTRoutes {
		if knownRoute.Match(method, path) {
			return true
		}
	}
	return false
}
