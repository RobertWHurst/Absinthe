package absinthe

import (
	"net/http"

	"github.com/nats-io/go-nats"
)

// Client manages the connection to Nats, as well as provides methods for
// binding routes and handlers, and dispatching HTTP requests and RPC calls
type Client struct {
	options  Options
	natsConn *nats.EncodedConn
	indexer  indexer
}

// Connect creates a new Client using the given Nats url, attempts to make a
// connection, then returns the connected client
func Connect(name, url string) *Client {
	options := GetDefaultOptions()
	options.Name = name
	options.URL = url

	c := &Client{options: options}
	c.connect()
	return c
}

// Call calls a RPC handler with the given name and arguments
func (c *Client) Call(name string, in interface{}, out interface{}) error {
	if err := c.indexer.validateCall(name, in, out); err != nil {
		return err
	}
	// TODO: Call the handler
}

func (c *Client) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

}

func (c *Client) subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return c.natsConn.Subscribe(c.options.Namespace+subj, cb)
}

func (c *Client) publish(subj string, v interface{}) error {
	return c.natsConn.Publish(c.options.Namespace+subj, v)
}

func (c *Client) connect() error {
	nc, err := c.options.getNatsOptions().Connect()
	if err != nil {
		return err
	}

	enc, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		return err
	}

	c.natsConn = enc
	c.indexer = indexer{client: c}
	go c.indexer.listenForNewHandlersAndRoutes()

	return nil
}
