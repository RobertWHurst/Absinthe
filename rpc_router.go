package absinthe

type RPCRouter struct {
	client *Client
}

func NewRPCRouter() *RPCRouter {
	return &RPCRouter{}
}
