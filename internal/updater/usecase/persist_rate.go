package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/updater/domain"
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
	snapshot domain.PriceSnapshot,
) error {
	return p.repo.Save(ctx, snapshot)
}
