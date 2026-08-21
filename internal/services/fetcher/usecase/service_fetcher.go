package usecase

import (
	"context"
	"errors"
	"log"
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
		return nil, errors.New("repo cannot be nil")
	}

	if client == nil {
		return nil, errors.New("client cannot be nil")
	}

	if producer == nil {
		return nil, errors.New("publisher cannot be nil")
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
		return nil, err
	}

	ctxWithTimeoutFetch, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	fetchedRates, err := f.client.Fetch(ctxWithTimeoutFetch, currencyList)
	if err != nil {
		return nil, err
	}

	err = f.producer.Produce(ctx, fetchedRates)
	if err != nil {
		return fetchedRates, err
	}

	return fetchedRates, nil
}

func (f *FetchService) RunScheduler(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_, err := f.UpdateRates(ctx)
			if err != nil {
				log.Printf("failed to update rates: %v", err)
				continue
			}
		}
	}
}

func (f *FetchService) FetchNewRates(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error) {
	ctxWithTimeoutFetch, cancelFetch := context.WithTimeout(ctx, 1*time.Second)
	defer cancelFetch()

	fetchedRates, err := f.client.Fetch(ctxWithTimeoutFetch, currencies)
	if err != nil {
		return nil, err
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
		resultErr = errors.Join(resultErr, err)
	}
	// 2. вызвать f.producer.Produce(ctx, rates), ошибку тоже сохранить
	ctxWithTimeoutProduce, cancelProduce := context.WithTimeout(ctx, 1*time.Second)
	defer cancelProduce()

	err = f.producer.Produce(ctxWithTimeoutProduce, fetchedRates)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	return fetchedRates, resultErr
}
