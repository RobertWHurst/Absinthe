package absinthe

import (
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const DefaultIndexerAnnouceInterval = 1 * time.Second

type indexer struct {
	client      *Client
	isAnouncing bool
	stoppedChan chan struct{}
	descriptor  *clientDescriptor
	peers       map[string]clientDescriptor
}

func createIndexer(client *Client) indexer {
	return indexer{
		client:      client,
		isAnouncing: false,
		stoppedChan: make(chan struct{}),
		descriptor:  &client.descriptor,
		peers:       make(map[string]clientDescriptor),
	}
}

func (i *indexer) start() {
	i.isAnouncing = true
	i.client.subscribe("announce", func(peer clientDescriptor) {
		i.addPeer(peer)
	})
	i.client.subscribe("rescind", func(peer clientDescriptor) {
		i.removePeerByID(peer.ID)
	})
	go func() {
		for i.isAnouncing {
			if err := i.client.publish("announce", *i.descriptor); err != nil {
				panic(err)
			}
			<-time.After(DefaultIndexerAnnouceInterval)
		}
		if err := i.client.publish("rescind", *i.descriptor); err != nil {
			panic(err)
		}
		i.stoppedChan <- struct{}{}
	}()
}

func (i *indexer) stop() {
	i.isAnouncing = false
	<-i.stoppedChan
}

func (i *indexer) addPeer(peer clientDescriptor) {
	if peer.ID != i.descriptor.ID {
		i.peers[peer.ID] = peer
	}
	spew.Dump(i.peers)
}

func (i *indexer) removePeerByID(peerID string) {
	if peerID != i.descriptor.ID {
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
