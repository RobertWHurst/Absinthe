package absinthe

import (
	"github.com/nats-io/go-nats"
)

const DefaultAbsintheNamespace = "__ABSINTHE__"

type Conn struct {
	*nats.EncodedConn
	Namespace string
}

func NewConn(options *Options) (*Conn, error) {
	natsConn, err := options.getNatsOptions().Connect()
	if err != nil {
		return nil, err
	}

	necodedNatsConn, err := nats.NewEncodedConn(natsConn, nats.GOB_ENCODER)
	if err != nil {
		return nil, err
	}

	return &Conn{
		EncodedConn: necodedNatsConn,
		Namespace:   options.Namespace,
	}, nil
}

func (c *Conn) Subscribe(subj string, cb nats.Handler) (*nats.Subscription, error) {
	return c.EncodedConn.Subscribe(c.Namespace+"."+subj, cb)
}

func (c *Conn) QueueSubscribe(subj string, queue string, cb nats.Handler) (*nats.Subscription, error) {
	return c.EncodedConn.QueueSubscribe(c.Namespace+"."+subj, queue, cb)
}

// TODO: Shadow the remaining message related nats connection methods so
// namespaces are always used correctly

func (c *Conn) Publish(subj string, v interface{}) error {
	return c.EncodedConn.Publish(c.Namespace+"."+subj, v)
}
