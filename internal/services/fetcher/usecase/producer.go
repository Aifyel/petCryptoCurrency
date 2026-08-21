package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type Producer interface {
	Produce(ctx context.Context, rate []entities.CurrencyRate) error
}
