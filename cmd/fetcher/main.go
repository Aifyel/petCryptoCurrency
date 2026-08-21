package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Aifyel/petCryptoCurrency/internal/services/fetcher/usecase"
	"github.com/Aifyel/petCryptoCurrency/internal/stub"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	repo := stub.NewRepository([]string{"USD", "MEOW"})
	client := stub.NewClient()
	producer := stub.NewProducer()

	fetchService, err := usecase.NewFetchService(repo, client, producer)
	if err != nil {
		panic(err)
	}

	scheduler := usecase.NewScheduler(fetchService, repo)

	fmt.Println("scheduler started, press Ctrl+C to stop")
	err = scheduler.Execute(ctx)
	if err != nil {
		panic(err)
	}

	final, err := repo.GetCurrencies(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("final repo state:", final)
}
