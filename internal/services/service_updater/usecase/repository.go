package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type RateRepository interface {
	Save(
		ctx context.Context,
		snapshot []entities.CurrencyRate,
	) error
}
