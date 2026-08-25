package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type FetcherClient interface {
	Fetch(ctx context.Context, currency []string) ([]entities.CurrencyRate, error)
}
