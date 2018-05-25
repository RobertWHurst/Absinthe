package absinthe

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// DefaultIndexerAnnouceInterval is the interval in number of seconds between
// each announcement of this client
const DefaultIndexerAnnouceInterval = 1 * time.Second

// Indexer keeps track of all absinthe instances (using the same namespace) on
// nats. It also announces the presence of this absinthe instance periotically.
// The indexer is used to validate rpc and rest requests, as well as obtain a
// peer to send the request to.
type Indexer struct {
	client      *Client
	isRunning   bool
	stoppedChan chan struct{}

	knownPeers map[string]Peer

	nextKnownPeers map[string]Peer
}

func NewIndexer(client *Client) Indexer {
	return Indexer{
		client:         client,
		isRunning:      false,
		stoppedChan:    make(chan struct{}),
		knownPeers:     make(map[string]Peer),
		nextKnownPeers: make(map[string]Peer),
	}
}

func (i *Indexer) HasRPCHandlerFor(path string) bool {
	for _, peer := range i.knownPeers {
		if peer.HasRPCHandlerFor(path) {
			return true
		}
	}
	return false
}

func (i *Indexer) HasRestHandlerFor(method, path string) bool {
	for _, peer := range i.knownPeers {
		if peer.HasRestHandlerFor(method, path) {
			return true
		}
	}
	return false
}

func (i *Indexer) Start() {
	i.isRunning = true

	i.client.Subscribe("PING", func(requestingPeerID string) {
		if requestingPeerID != i.client.ID {
			err := i.client.Publish("PONG-"+requestingPeerID, i.client.Peer)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	i.client.Subscribe("PONG-"+i.client.ID, func(respondingPeer Peer) {
		if respondingPeer.ID != i.client.ID {
			i.knownPeers[respondingPeer.ID] = respondingPeer
			i.nextKnownPeers[respondingPeer.ID] = respondingPeer
		}
	})

	for i.isRunning {
		i.client.Publish("PING", i.client.ID)
		<-time.After(DefaultIndexerAnnouceInterval)
		for k := range i.knownPeers {
			delete(i.knownPeers, k)
		}
		for k := range i.nextKnownPeers {
			i.knownPeers[k] = i.nextKnownPeers[k]
		}
		for k := range i.nextKnownPeers {
			delete(i.nextKnownPeers, k)
		}
		spew.Dump(i.knownPeers)
	}
	i.stoppedChan <- struct{}{}
}

func (i *Indexer) Stop() {
	i.isRunning = false
	<-i.stoppedChan
}
