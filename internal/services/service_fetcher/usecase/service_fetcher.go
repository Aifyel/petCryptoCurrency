package usecase

import (
	"context"
	errorsBase "errors"
	"time"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(entities.ErrInvalidParams, "service_fetcher Repository:")
	}

	if client == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "service_fetcher Client:")
	}

	if producer == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "service_fetcher Producer:")
	}

	return &FetchService{
		repo:     repo,
		client:   client,
		producer: producer,
	}, nil
}

func (f *FetchService) UpdateRates(
	ctx context.Context,
) error {
	ctxWithTimeoutRepo, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	currencyList, err := f.repo.GetCurrencies(ctxWithTimeoutRepo)
	if err != nil {
		return errors.Wrap(entities.ErrRepositoryFailure, "service_fetcher Repository:")
	}

	ctxWithTimeoutFetch, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	fetchedRates, err := f.client.Fetch(ctxWithTimeoutFetch, currencyList)
	if err != nil {
		return errors.Wrap(entities.ErrClientFailure, "service_fetcher Client:")
	}

	err = f.producer.Produce(ctx, fetchedRates, TopicScheduled)
	if err != nil {
		return errors.Wrap(entities.ErrClientFailure, "service_fetcher Broker:")
	}

	return nil
}

func (f *FetchService) FetchNewRates(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error) {
	ctxWithTimeoutFetch, cancelFetch := context.WithTimeout(ctx, 1*time.Second)
	defer cancelFetch()

	fetchedRates, err := f.client.Fetch(ctxWithTimeoutFetch, currencies)
	if err != nil {
		return nil, errors.Wrap(entities.ErrClientFailure, "service_fetcher Client:")
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
		resultErr = errorsBase.Join(resultErr, errors.Wrap(entities.ErrRepositoryFailure, "service_fetcher SaveNewCurrency:"))
	}

	ctxWithTimeoutProduce, cancelProduce := context.WithTimeout(ctx, 1*time.Second)
	defer cancelProduce()

	err = f.producer.Produce(ctxWithTimeoutProduce, fetchedRates, TopicNew)
	if err != nil {
		resultErr = errorsBase.Join(resultErr, errors.Wrap(entities.ErrMessagingFailure, "service_fetcher Produce:"))
	}

	return fetchedRates, resultErr
}
