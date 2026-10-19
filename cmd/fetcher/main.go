package main

import (
	"context"
	"log"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	err = app.Run(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
