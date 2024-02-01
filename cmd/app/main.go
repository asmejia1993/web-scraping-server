package main

import (
	"context"
	"log"

	"github.com/asmejia1993/web-scraping-server/pkg/http/rest"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(ctx context.Context) error {
	server, err := rest.NewServer(ctx)
	if err != nil {
		return err
	}
	err = server.Run(ctx)
	return err
}
