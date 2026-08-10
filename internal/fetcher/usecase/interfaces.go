package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/fetcher/domain"
)

type MarketDataClient interface {
	Fetch(ctx context.Context, currencies []domain.Currency) ([]domain.CurrencyRate, error)
}

type CurrencyRepository interface {
	GetCurrencies(ctx context.Context) ([]domain.Currency, error)
	SaveNewCurrency(ctx context.Context, currency domain.Currency) error
	Exists(ctx context.Context, currency domain.Currency) (bool, error)
}

type RatePublisher interface {
	Publish(ctx context.Context, rate domain.CurrencyRate) error
}

type FetchCurrencyGRPC interface {
	Execute(ctx context.Context, currency domain.Currency) (domain.CurrencyRate, error)
}

type FetchMarketRatesUseCase interface {
	Execute(ctx context.Context) error
}
