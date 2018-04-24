package absinthe

import (
	"reflect"
	"time"
)

const indexerAnnouceInterval = 2 * time.Minute

type indexer struct {
	client      *Client
	isListening bool
	descriptor  *clientDescriptor
	peers       map[string]clientDescriptor
}

func (i *indexer) listenForNewHandlersAndRoutes() {
	i.client.subscribe("anounce", func(peer clientDescriptor) { i.addPeer(peer) })
	i.client.subscribe("rescind", func(peer clientDescriptor) { i.removePeerByID(peer.id) })
	go func() {
		for i.isListening {
			i.client.publish("anounce", i.descriptor)
			<-time.After(indexerAnnouceInterval)
		}
		i.client.publish("rescind", i.descriptor)
	}()
}

func (i *indexer) addPeer(peer clientDescriptor) {
	if peer.id != i.descriptor.id {
		i.peers[peer.id] = peer
	}
}

func (i *indexer) removePeerByID(peerID string) {
	if peerID != i.descriptor.id {
		delete(i.peers, peerID)
	}
}

func (i *indexer) validateCall(namespace string, inType, outType reflect.Type) bool {
	if i.descriptor.testCall(namespace, inType, outType) {
		return true
	}
	for _, peer := range i.peers {
		if peer.testCall(namespace, inType, outType) {
			return true
		}
	}
	return false
}

func (i *indexer) validateRequest() bool {
	if i.descriptor.testRequest() {
		return true
	}
	for _, peer := range i.peers {
		if peer.testRequest() {
			return true
		}
	}
	return false
}

func (i *indexer) close() {
	i.isListening = false
}
