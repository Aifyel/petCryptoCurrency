package usecase

import (
	"context"
	"fmt"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type GetCurrencyStatistics struct {
	repo   Repository
	client FetcherClient
}

func NewGetCurrencyStatistics(
	repo Repository,
	client FetcherClient,
) (*GetCurrencyStatistics, error) {
	if client == nil {
		return nil, fmt.Errorf("server client: %w", entities.ErrInvalidParams)
	}

	if repo == nil {
		return nil, fmt.Errorf("server repository: %w", entities.ErrInvalidParams)
	}

	return &GetCurrencyStatistics{
		repo:   repo,
		client: client,
	}, nil
}

func (g *GetCurrencyStatistics) GetStatistics(
	ctx context.Context,
	currency []string,
) ([]entities.CurrencyStatistics, error) {
	if len(currency) == 0 {
		return nil, fmt.Errorf("server empty currency slice: %w", entities.ErrInvalidParams)
	}

	stats, err := g.repo.GetStatistics(ctx, currency)
	if err != nil {
		return nil, fmt.Errorf("server GetStatistics: %w, %w", entities.ErrRepositoryFailure, err)
	}

	currencyMap := make(map[string]struct{})
	for _, rate := range stats {
		currencyMap[rate.Currency] = struct{}{}
	}

	missingArr := make([]string, 0, len(currency))
	for _, c := range currency {
		if _, ok := currencyMap[c]; !ok {
			missingArr = append(missingArr, c)
		}
	}

	if len(missingArr) > 0 {
		newRates := make([]entities.CurrencyRate, 0, len(missingArr))
		if len(missingArr) > 0 {
			rates, err := g.client.Fetch(ctx, missingArr)
			if err != nil {
				return nil, fmt.Errorf("server Fetch: %w, %w", entities.ErrClientFailure, err)
			}
			newRates = append(newRates, rates...)
		}

		for _, rate := range newRates {
			stats = append(stats, entities.CurrencyStatistics{
				Currency:          rate.Currency,
				CurrentPrice:      rate.Price,
				LowestPriceDaily:  rate.Price,
				HighestPriceDaily: rate.Price,
				ChangeHourPercent: 0,
			})
		}
	}

	return stats, nil
}
