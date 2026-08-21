package usecase

import (
	"context"

	fetcher "github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/server_api/domain"
)

type Repository interface {
	GetStatistics(
		ctx context.Context,
		currency domain.Currency,
	) (fetcher.CurrencyStatistics, error)
}
