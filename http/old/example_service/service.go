package main

import (
	"fmt"

	"github.com/RobertWHurst/absinthe"
	"github.com/RobertWHurst/absinthe/abshttp"
	"github.com/nats-io/go-nats"
)

// import (
// 	absinthe "github.com/RobertWHurst/Absinthe"
// 	nats "github.com/nats-io/go-nats"
// )

func main() {
	nc, err := nats.GetDefaultOptions().Connect()
	if err != nil {
		panic(err)
	}

	abs, err := abs.FromNatsConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	abs.GET("/", func(c *abshttp.Context) {
		fmt.Printf("%+v", c)
		c.Response.StatusCode = 200
		c.Body = "HELLO FROM THE SERVICE!!!"
		c.End()
	})

	select {}
}
