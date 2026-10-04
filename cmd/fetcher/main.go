package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/adapters/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("%v", err)
	}

	conf, err := config.Load()
	if err != nil {
		log.Fatalf("%v", err)
	}

	pool, err := pgxpool.New(ctx, conf.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create pgxpool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	repo, err := postgres.NewRepository(pool)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	err = repo.SaveNewCurrency(ctx, []string{"bitcoin", "ethereum"})
	if err != nil {
		log.Fatalf("failed to save new currency: %v", err)
	}

	currencies, err := repo.GetCurrencies(ctx)
	if err != nil {
		log.Fatalf("failed to fetch currencies: %v", err)
	}

	fmt.Printf("Currencies: %v\n", currencies)
}
