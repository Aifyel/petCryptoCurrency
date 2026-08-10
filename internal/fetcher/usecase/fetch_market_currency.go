package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/fetcher/domain"
)

type FetchCurrency struct {
	repo      CurrencyRepository
	client    MarketDataClient
	publisher RatePublisher
}

func NewFetchCurrency(
	repo CurrencyRepository,
	client MarketDataClient,
	publisher RatePublisher,
) *FetchCurrency {
	return &FetchCurrency{
		repo:      repo,
		client:    client,
		publisher: publisher,
	}
}

func (f *FetchCurrency) Execute(
	ctx context.Context,
	currency domain.Currency,
) (domain.CurrencyRate, error) {
	if currency == "" {
		//err
	}

	exists, err := f.repo.Exists(ctx, currency)
	if err != nil {
		return domain.CurrencyRate{}, err
	}

	if exists {
		return domain.CurrencyRate{}, nil //err
	}

	return domain.CurrencyRate{} /*nil?*/, nil //err
}
