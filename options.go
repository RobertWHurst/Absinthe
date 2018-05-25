package absinthe

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/coreos/go-semver/semver"

	nats "github.com/nats-io/go-nats"
)

// Options can be used to configure absinthe clients
type Options struct {
	// - Absinthe Options -

	// IndexingInterval controls how often absinthe will index handlers and
	// routes.
	IndexingInterval time.Duration

	// Name is an optional identifier for the absinthe client used to identify
	// it on the Nats server
	Name string

	// Version is an optional identifier for the absinthe client used to identify
	// it on the Nats server
	Version *semver.Version

	// Namespace is used as a prefix for Nats subjects to scope absinthe specific
	// messages
	Namespace string

	// - Nats Options -

	// Servers is a configured set of servers which this client
	// will use when attempting to connect.
	Servers []string

	// NoRandomize configures whether we will randomize the
	// server pool.
	NoRandomize bool

	// Verbose signals the server to send an OK ack for commands
	// successfully processed by the server.
	Verbose bool

	// Pedantic signals the server whether it should be doing further
	// validation of subjects.
	Pedantic bool

	// Secure enables TLS secure connections that skip server
	// verification by default. NOT RECOMMENDED.
	Secure bool

	// TLSConfig is a custom TLS configuration to use for secure
	// transports.
	TLSConfig *tls.Config

	// AllowReconnect enables reconnection logic to be used when we
	// encounter a disconnect from the current server.
	AllowReconnect bool

	// MaxReconnect sets the number of reconnect attempts that will be
	// tried before giving up. If negative, then it will never give up
	// trying to reconnect.
	MaxReconnect int

	// ReconnectWait sets the time to backoff after attempting a reconnect
	// to a server that we were already connected to previously.
	ReconnectWait time.Duration

	// Timeout sets the timeout for a Dial operation on a connection.
	Timeout time.Duration

	// FlusherTimeout is the maximum time to wait for the flusher loop
	// to be able to finish writing to the underlying connection.
	FlusherTimeout time.Duration

	// PingInterval is the period at which the client will be sending ping
	// commands to the server, disabled if 0 or negative.
	PingInterval time.Duration

	// MaxPingsOut is the maximum number of pending ping commands that can
	// be awaiting a response before raising an ErrStaleConnection error.
	MaxPingsOut int

	// ClosedCB sets the closed handler that is called when a client will
	// no longer be connected.
	ClosedCB nats.ConnHandler

	// DisconnectedCB sets the disconnected handler that is called
	// whenever the connection is disconnected.
	DisconnectedCB nats.ConnHandler

	// ReconnectedCB sets the reconnected handler called whenever
	// the connection is successfully reconnected.
	ReconnectedCB nats.ConnHandler

	// DiscoveredServersCB sets the callback that is invoked whenever a new
	// server has joined the cluster.
	DiscoveredServersCB nats.ConnHandler

	// AsyncErrorCB sets the async error handler (e.g. slow consumer errors)
	AsyncErrorCB nats.ErrHandler

	// ReconnectBufSize is the size of the backing bufio during reconnect.
	// Once this has been exhausted publish operations will return an error.
	ReconnectBufSize int

	// SubChanLen is the size of the buffered channel used between the socket
	// Go routine and the message delivery for SyncSubscriptions.
	// NOTE: This does not affect AsyncSubscriptions which are
	// dictated by PendingLimits()
	SubChanLen int

	// User sets the username to be used when connecting to the server.
	User string

	// Password sets the password to be used when connecting to a server.
	Password string

	// Token sets the token to be used when connecting to a server.
	Token string

	// CustomDialer allows to specify a custom dialer (not necessarily
	// a *net.Dialer).
	CustomDialer nats.CustomDialer
}

type Option func(*Options) error

// GetDefaultOptions returns the default options for absinthe clients
func GetDefaultOptions() Options {
	natsOptionDefaults := nats.GetDefaultOptions()
	interval, _ := time.ParseDuration("5s")

	return Options{
		IndexingInterval: interval,
		Namespace:        "absinthe",

		Servers:             natsOptionDefaults.Servers,
		NoRandomize:         natsOptionDefaults.NoRandomize,
		Verbose:             natsOptionDefaults.Verbose,
		Pedantic:            natsOptionDefaults.Pedantic,
		Secure:              natsOptionDefaults.Secure,
		TLSConfig:           natsOptionDefaults.TLSConfig,
		AllowReconnect:      natsOptionDefaults.AllowReconnect,
		MaxReconnect:        natsOptionDefaults.MaxReconnect,
		ReconnectWait:       natsOptionDefaults.ReconnectWait,
		Timeout:             natsOptionDefaults.Timeout,
		FlusherTimeout:      natsOptionDefaults.FlusherTimeout,
		PingInterval:        natsOptionDefaults.PingInterval,
		MaxPingsOut:         natsOptionDefaults.MaxPingsOut,
		ClosedCB:            natsOptionDefaults.ClosedCB,
		DisconnectedCB:      natsOptionDefaults.DisconnectedCB,
		ReconnectedCB:       natsOptionDefaults.ReconnectedCB,
		DiscoveredServersCB: natsOptionDefaults.DiscoveredServersCB,
		AsyncErrorCB:        natsOptionDefaults.AsyncErrorCB,
		ReconnectBufSize:    natsOptionDefaults.ReconnectBufSize,
		SubChanLen:          natsOptionDefaults.SubChanLen,
		User:                natsOptionDefaults.User,
		Password:            natsOptionDefaults.Password,
		Token:               natsOptionDefaults.Token,
		CustomDialer:        natsOptionDefaults.CustomDialer,
	}
}

// Connect creates a client and connects to nats using the options it is called
// upon. It then returns the connected client.
func (o Options) Connect() (*Client, error) {
	c := &Client{options: o}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (o Options) getNatsOptions() nats.Options {
	natsOptions := nats.GetDefaultOptions()

	natsOptions.Servers = o.Servers
	natsOptions.NoRandomize = o.NoRandomize
	natsOptions.Name = o.Name
	natsOptions.Verbose = o.Verbose
	natsOptions.Pedantic = o.Pedantic
	natsOptions.Secure = o.Secure
	natsOptions.TLSConfig = o.TLSConfig
	natsOptions.AllowReconnect = o.AllowReconnect
	natsOptions.MaxReconnect = o.MaxReconnect
	natsOptions.ReconnectWait = o.ReconnectWait
	natsOptions.Timeout = o.Timeout
	natsOptions.FlusherTimeout = o.FlusherTimeout
	natsOptions.PingInterval = o.PingInterval
	natsOptions.MaxPingsOut = o.MaxPingsOut
	natsOptions.ClosedCB = o.ClosedCB
	natsOptions.DisconnectedCB = o.DisconnectedCB
	natsOptions.ReconnectedCB = o.ReconnectedCB
	natsOptions.DiscoveredServersCB = o.DiscoveredServersCB
	natsOptions.AsyncErrorCB = o.AsyncErrorCB
	natsOptions.ReconnectBufSize = o.ReconnectBufSize
	natsOptions.SubChanLen = o.SubChanLen
	natsOptions.User = o.User
	natsOptions.Password = o.Password
	natsOptions.Token = o.Token
	natsOptions.CustomDialer = o.CustomDialer

	return natsOptions
}

func Name(name string) Option {
	return func(o *Options) error {
		o.Name = name
		return nil
	}
}

func Version(version string) Option {
	return func(o *Options) error {
		version, err := semver.NewVersion(version)
		o.Version = version
		return err
	}
}

func Namespace(namespace string) Option {
	return func(o *Options) error {
		o.Namespace = namespace
		return nil
	}
}

func DontRandomize() Option {
	return func(o *Options) error {
		o.NoRandomize = true
		return nil
	}
}

func Secure(tls ...*tls.Config) Option {
	return func(o *Options) error {
		o.Secure = true
		if len(tls) > 1 {
			return nats.ErrMultipleTLSConfigs
		}
		if len(tls) == 1 {
			o.TLSConfig = tls[0]
		}
		return nil
	}
}

func RootCAs(file ...string) Option {
	return func(o *Options) error {
		pool := x509.NewCertPool()
		for _, f := range file {
			rootPEM, err := ioutil.ReadFile(f)
			if err != nil || rootPEM == nil {
				return fmt.Errorf("absinthe: error loading or parsing rootCA file: %v", err)
			}
			ok := pool.AppendCertsFromPEM(rootPEM)
			if !ok {
				return fmt.Errorf("absinthe: failed to parse root certificate from %q", f)
			}
		}
		if o.TLSConfig == nil {
			o.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		o.TLSConfig.RootCAs = pool
		o.Secure = true
		return nil
	}
}

func ClientCert(certFile, keyFile string) Option {
	return func(o *Options) error {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("absinthe: error loading client certificate: %v", err)
		}
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return fmt.Errorf("absinthe: error parsing client certificate: %v", err)
		}
		if o.TLSConfig == nil {
			o.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		o.TLSConfig.Certificates = []tls.Certificate{cert}
		o.Secure = true
		return nil
	}
}

func NoReconnect() Option {
	return func(o *Options) error {
		o.AllowReconnect = false
		return nil
	}
}

func ReconnectWait(t time.Duration) Option {
	return func(o *Options) error {
		o.ReconnectWait = t
		return nil
	}
}

func MaxReconnects(max int) Option {
	return func(o *Options) error {
		o.MaxReconnect = max
		return nil
	}
}

func PingInterval(t time.Duration) Option {
	return func(o *Options) error {
		o.PingInterval = t
		return nil
	}
}

func ReconnectBufSize(size int) Option {
	return func(o *Options) error {
		o.ReconnectBufSize = size
		return nil
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) error {
		o.Timeout = t
		return nil
	}
}

func DisconnectHandler(cb nats.ConnHandler) Option {
	return func(o *Options) error {
		o.DisconnectedCB = cb
		return nil
	}
}

func ReconnectHandler(cb nats.ConnHandler) Option {
	return func(o *Options) error {
		o.ReconnectedCB = cb
		return nil
	}
}

func ClosedHandler(cb nats.ConnHandler) Option {
	return func(o *Options) error {
		o.ClosedCB = cb
		return nil
	}
}

// DiscoveredServersHandler is an Option to set the new servers handler.
func DiscoveredServersHandler(cb nats.ConnHandler) Option {
	return func(o *Options) error {
		o.DiscoveredServersCB = cb
		return nil
	}
}

func UserInfo(user, password string) Option {
	return func(o *Options) error {
		o.User = user
		o.Password = password
		return nil
	}
}

func Token(token string) Option {
	return func(o *Options) error {
		o.Token = token
		return nil
	}
}

func SetCustomDialer(dialer nats.CustomDialer) Option {
	return func(o *Options) error {
		o.CustomDialer = dialer
		return nil
	}
}

const DefaultURL = nats.DefaultURL
