package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type PersistRate struct {
	repo RateRepository
}

func NewPersistRate(repo RateRepository) *PersistRate {
	return &PersistRate{
		repo: repo,
	}
}

func (p *PersistRate) Execute(
	ctx context.Context,
	snapshot entities.CurrencyRate,
) error {
	return p.repo.Save(ctx, snapshot)
}
