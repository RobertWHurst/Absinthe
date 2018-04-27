package absinthe

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/nats-io/go-nats"
)

// Client manages the connection to Nats, as well as provides methods for
// binding routes and handlers, and dispatching HTTP requests and RPC calls
type Client struct {
	options    Options
	natsConn   *nats.EncodedConn
	indexer    indexer
	descriptor clientDescriptor
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
	options.Servers = processUrlString(url)

	return options.Connect()
}

// Call calls a RPC handler with the given name and arguments
func (c *Client) Call(name string, in interface{}, out interface{}) error {
	inType := reflect.TypeOf(in)
	outType := reflect.TypeOf(out)
	if c.indexer.validateCall(name, inType, outType) {

	}
	// TODO: Call the handler
	return nil
}

func (c *Client) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	spew.Dump("%+v", c)
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

	enc.Subscribe("*", func(v *interface{}) {
		fmt.Println("Message")
		spew.Dump(v)
	})
	c.natsConn = enc
	c.indexer = createIndexer(c)
	go c.indexer.start()

	return nil
}

func (c *Client) subscribe(subj string, cb nats.Handler) (*nats.Subscription, error) {
	return c.natsConn.Subscribe(c.getNamespace()+"."+subj, cb)
}

func (c *Client) publish(subj string, v interface{}) error {
	return c.natsConn.Publish(c.getNamespace()+"."+subj, v)
}

func (c *Client) flush() error {
	return c.natsConn.Flush()
}

func (c *Client) getNamespace() string {
	if len(c.options.Namespace) != 0 {
		return c.options.Namespace
	}
	return "__ABSINTHE__"
}

func processUrlString(url string) []string {
	urls := strings.Split(url, ",")
	for i, s := range urls {
		urls[i] = strings.TrimSpace(s)
	}
	return urls
}
