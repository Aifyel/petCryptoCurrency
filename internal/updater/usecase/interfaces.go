package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/updater/domain"
)

type RateRepository interface {
	Save(
		ctx context.Context,
		snapshot domain.PriceSnapshot,
	) error
}
