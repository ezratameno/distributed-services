package main

import (
	"fmt"
	"os"

	"github.com/ezratameno/distributed-services/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	srv := server.NewHTTPServer(":8080")
	return srv.ListenAndServe()
}
