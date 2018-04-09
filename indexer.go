package absinthe

type indexer struct {
	client         *Client
	isListening    bool
	remoteHandlers []remoteHandler
	remoteRoutes   []remoteRoute
}

func (i *indexer) listenForNewHandlersAndRoutes() {
	i.client.subscribe("")
	i.client.publish("joined")
}

func (i *indexer) close() {
	i.isListening = false
}
