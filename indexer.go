package absinthe

type indexer struct {
	client      *Client
	isListening bool
	peers       []clientDescriptor
}

func (i *indexer) listenForNewHandlersAndRoutes() {
	i.client.subscribe("")
	i.client.publish("joined")
}

func (i *indexer) close() {
	i.isListening = false
}
