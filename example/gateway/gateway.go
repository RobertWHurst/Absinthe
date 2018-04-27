package main

import (
	"net/http"

	"github.com/RobertWHurst/Absinthe"
)

func main() {
	client, err := absinthe.Connect(
		absinthe.DefaultURL,
		absinthe.Name("gateway"),
		absinthe.Version("0.1.0"),
	)

	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:    ":8000",
		Handler: client,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
