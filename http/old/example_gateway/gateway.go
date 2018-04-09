package main

import (
	"fmt"
	"net/http"

	abs "github.com/RobertWHurst/absinthe"
	"github.com/nats-io/go-nats"
)

func main() {
	fmt.Print("Connecting to Nats")
	nc, err := nats.GetDefaultOptions().Connect()
	if err != nil {
		panic(err)
	}
	fmt.Print("Connected to Nats")

	fmt.Print("Creating Absinthe connection")
	abs, _ := abs.FromNatsConn(nc, nats.JSON_ENCODER)
	fmt.Print("Connection created")

	fmt.Print("Creating HTTP server")
	server := http.Server{
		Addr:    ":8080",
		Handler: abs,
	}
	fmt.Print("HTTP server created")

	fmt.Print("Binding a socket")
	server.ListenAndServe()
	fmt.Print("Socket bound")
}
