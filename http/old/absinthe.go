package abs

// TODO: Add config value for time outs where HTTP timeout and Nats subscription
// timeout are set from. HTTP timeout Nats subscription timeout should be be
// same.

import (
	"io"
	"net/http"
	"strings"

	"github.com/RobertWHurst/absinthe/abshttp"
	"github.com/nats-io/go-nats"
)

const BODY_CHUNK_SIZE = 1024 * 32
const BODY_BUFFER_SIZE = BODY_CHUNK_SIZE * 10

// FromNatsConn creates a nats.Conn from a nats Conn instance
func FromNatsConn(nc *nats.Conn, encodeType string) (*Conn, error) {
	encodedNc, err := nats.NewEncodedConn(nc, encodeType)
	if err != nil {
		return &Conn{}, err
	}
	return FromEncodedNatsConn(encodedNc), nil
}

// FromEncodedNatsConn creates an nats.Conn from a nats EncodedConn instance
func FromEncodedNatsConn(encodedNc *nats.EncodedConn) *Conn {
	conn := &Conn{
		Router: Router{
			root:     true,
			basePath: "",
		},
		natsConn:        encodedNc.Conn,
		natsEncodedConn: encodedNc,
	}
	conn.Router.conn = conn

	return conn
}

// Conn represents the Absinthe network of services. It has a collection of
// methods for binding gin routes, handling http requests, and
// binding/dispatching RPC.
type Conn struct {
	Router
	natsConn        *nats.Conn
	natsEncodedConn *nats.EncodedConn
}

func (c *Conn) Call(rpcHandlerName string, args *interface{}, reply *interface{}) {

}

func (c *Conn) ServeHTTP(originalResponseWriter http.ResponseWriter, originalRequest *http.Request) {
	context := abshttp.ContextFromHTTP(originalResponseWriter, originalRequest)

	// ----
	// ----
	// ----
	// ----
	// ----

	nc := c.natsEncodedConn

	// Create a Absinthe HTTP request object. This object will be sent over
	// Nats to any services that might be able to respond to it.
	request := abshttp.NewRequest(originalRequest)

	// Listen for the response and collect it
	responseChan := make(chan *abshttp.Response)
	responseBodyChan := make(chan []byte)
	nc.Subscribe(request.Subjects.Response, func(response *abshttp.Response) {
		responseChan <- response
	})
	nc.BindRecvChan(request.Subjects.ResponseBody, responseBodyChan)

	// Send the HTTP request into Nats
	nc.Publish(request.Subjects.Request, request)
	bodyChan := make(chan []byte, BODY_BUFFER_SIZE)
	nc.BindSendChan(request.Subjects.RequestBody, bodyChan)
	writeReaderToByteChan(originalRequest.Body, bodyChan)

	// Block until we have a response
	response := <-responseChan
	originalResponseWriter.WriteHeader(response.StatusCode)

	// Write out the response body to the client
	for {
		chunk, ok := <-responseBodyChan
		println("Chunk", chunk, ok)
		if len(chunk) == 0 || !ok {
			println("Done")
			break
		}
		originalResponseWriter.Write(chunk)
	}
}

func (c *Conn) httpSubscribe(method string, absolutePath string, handlers abshttp.HandlerChain) {
	nc := c.natsEncodedConn
	subject := abshttp.NewRequestSubject(method, absolutePath)

	// Listen for incoming HTTP requests for the route we are binding
	nc.Subscribe(subject, func(request *abshttp.Request) {

		// Create a context for all handlers of the route to use in order to respond
		// to the request.
		context := abshttp.NewContext(request)

		// Execute the handler tree
		// TODO: Implement the handler tree
		handlers[0](context)

		// Send the response to the absinthe gateway
		nc.Publish(context.Request.Subjects.Response, context.Response)

		// Send the response body to the absinthe gateway
		bodyChan := make(chan []byte, BODY_BUFFER_SIZE)
		nc.BindSendChan(request.Subjects.ResponseBody, bodyChan)
		writeReaderToByteChan(strings.NewReader(context.Body), bodyChan)
	})
}

func writeReaderToByteChan(bodyReader io.Reader, bodyChan chan []byte) error {
	for {
		bodyChunk := make([]byte, BODY_CHUNK_SIZE)
		n, err := bodyReader.Read(bodyChunk)
		if err != nil {
			// TODO: raise errors
			break
		}
		bodyChan <- bodyChunk
		if n != BODY_CHUNK_SIZE {
			close(bodyChan)
			break
		}
	}
	return nil
}
