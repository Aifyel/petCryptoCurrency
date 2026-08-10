package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/reader/domain"
)

type GetCurrencyStatistics struct {
	repo    RateQueryRepository
	fetcher FetcherClient
}

func NewGetCurrencyStatistics(
	repo RateQueryRepository,
	fetcher FetcherClient,
) *GetCurrencyStatistics {
	return &GetCurrencyStatistics{
		repo:    repo,
		fetcher: fetcher,
	}
}

func (g *GetCurrencyStatistics) GetStatistics(
	ctx context.Context,
	currency domain.Currency,
) (domain.CurrencyStatistics, error) {
	return domain.CurrencyStatistics{}, nil
}
