package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type Repository interface {
	GetStatistics(ctx context.Context, currency []string) ([]entities.CurrencyStatistics, error)
}
