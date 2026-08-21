package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type APIClient interface {
	Fetch(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error)
}
