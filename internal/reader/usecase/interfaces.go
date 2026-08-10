package usecase

import (
	"context"

	fetcher "github.com/Aifyel/petCryptoCurrency/internal/fetcher/domain"
	"github.com/Aifyel/petCryptoCurrency/internal/reader/domain"
)

type RateQuery interface {
	GetStatistics(
		ctx context.Context,
		currency domain.Currency,
	) (domain.CurrencyStatistics, error)
}

type RateQueryRepository interface {
	GetStatistics(
		ctx context.Context,
		currency domain.Currency,
	) (domain.CurrencyStatistics, error)
}

type FetcherClient interface {
	FetchCurrency(
		ctx context.Context,
		currency domain.Currency,
	) (fetcher.CurrencyRate, error)
}
