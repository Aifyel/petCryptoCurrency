package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type FetchService struct {
	repo     Repository
	client   APIClient
	producer Producer
}

func NewFetchService(
	repo Repository,
	client APIClient,
	producer Producer,
) (*FetchService, error) {
	if repo == nil {
		return nil, fmt.Errorf("fetcher repository: %w", entities.ErrInvalidParams)
	}

	if client == nil {
		return nil, fmt.Errorf("fetcher client: %w", entities.ErrInvalidParams)
	}

	if producer == nil {
		return nil, fmt.Errorf("fetcher producer: %w", entities.ErrInvalidParams)
	}

	return &FetchService{
		repo:     repo,
		client:   client,
		producer: producer,
	}, nil
}

func (f *FetchService) UpdateRates(
	ctx context.Context,
) ([]entities.CurrencyRate, error) {
	ctxWithTimeoutRepo, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	currencyList, err := f.repo.GetCurrencies(ctxWithTimeoutRepo)
	if err != nil {
		return nil, fmt.Errorf("fetcher GetCurrencies: %w, %w", entities.ErrRepositoryFailure, err)
	}

	ctxWithTimeoutFetch, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	fetchedRates, err := f.client.Fetch(ctxWithTimeoutFetch, currencyList)
	if err != nil {
		return nil, fmt.Errorf("fetcher Fetch: %w, %w", entities.ErrClientFailure, err)
	}

	err = f.producer.Produce(ctx, fetchedRates)
	if err != nil {
		return fetchedRates, fmt.Errorf("fetcher Produce: %w, %w", entities.ErrMessagingFailure, err)
	}

	return fetchedRates, nil
}

func (f *FetchService) FetchNewRates(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error) {
	ctxWithTimeoutFetch, cancelFetch := context.WithTimeout(ctx, 1*time.Second)
	defer cancelFetch()

	fetchedRates, err := f.client.Fetch(ctxWithTimeoutFetch, currencies)
	if err != nil {
		return nil, fmt.Errorf("fetcher Fetch: %w, %w", entities.ErrClientFailure, err)
	}

	currency := make([]string, 0, len(fetchedRates))
	for _, rate := range fetchedRates {
		currency = append(currency, rate.Currency)
	}

	ctxWithTimeoutRepo, cancelRepo := context.WithTimeout(ctx, 1*time.Second)
	defer cancelRepo()

	var resultErr error

	err = f.repo.SaveNewCurrency(ctxWithTimeoutRepo, currency)
	if err != nil {
		resultErr = errors.Join(resultErr, fmt.Errorf("fetcher SaveNewCurrency: %w, %w", entities.ErrRepositoryFailure, err))
	}

	ctxWithTimeoutProduce, cancelProduce := context.WithTimeout(ctx, 1*time.Second)
	defer cancelProduce()

	err = f.producer.Produce(ctxWithTimeoutProduce, fetchedRates)
	if err != nil {
		resultErr = errors.Join(resultErr, fmt.Errorf("fetcher Produce: %w, %w", entities.ErrMessagingFailure, err))
	}

	return fetchedRates, resultErr
}
