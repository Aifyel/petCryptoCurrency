package usecase

import (
	"context"
)

type FetchMarketRates struct {
	repo      CurrencyRepository
	client    MarketDataClient
	publisher RatePublisher
}

func NewFetchMarketRates(
	repo CurrencyRepository,
	client MarketDataClient,
	publisher RatePublisher,
) *FetchMarketRates {
	return &FetchMarketRates{
		repo:      repo,
		client:    client,
		publisher: publisher,
	}
}

func (f *FetchMarketRates) Execute(ctx context.Context) error {
	return nil
}
