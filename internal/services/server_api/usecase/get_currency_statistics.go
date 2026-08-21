package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/server_api/domain"
)

type GetCurrencyStatistics struct {
	repo Repository
}

func NewGetCurrencyStatistics(
	repo Repository,
) *GetCurrencyStatistics {
	return &GetCurrencyStatistics{
		repo: repo,
	}
}

func (g *GetCurrencyStatistics) GetStatistics(
	ctx context.Context,
	currency domain.Currency,
) (entities.CurrencyStatistics, error) {
	return entities.CurrencyStatistics{}, nil
}
