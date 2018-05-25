package absinthe

import (
	"strings"
)

// Client manages the connection to Nats, as well as provides methods for
// binding routes and handlers, and dispatching HTTP requests and RPC calls
type Client struct {
	Peer
	Conn
	*RESTRouter
	*RPCRouter
	options Options
	indexer Indexer
}

// Connect creates a new Client using the given Nats url, attempts to make a
// connection, then returns the connected client
func Connect(url string, optionSetters ...Option) (*Client, error) {
	options := GetDefaultOptions()

	for _, optionSetter := range optionSetters {
		if err := optionSetter(&options); err != nil {
			return nil, err
		}
	}
	options.Servers = processURLString(url)

	return options.Connect()
}

func (c *Client) connect() error {
	peer, _ := NewPeer(c.options.Name, c.options.Version, "")
	c.Peer = *peer

	conn, err := NewConn(&c.options)
	if err != nil {
		return err
	}
	c.Conn = *conn

	c.indexer = NewIndexer(c)
	c.RESTRouter = NewRESTRouter()
	c.RESTRouter.client = c
	c.RPCRouter = NewRPCRouter()
	c.RPCRouter.client = c

	go c.indexer.Start()

	return nil
}

func processURLString(url string) []string {
	urls := strings.Split(url, ",")
	for i, s := range urls {
		urls[i] = strings.TrimSpace(s)
	}
	return urls
}
