package postgres

import (
	"context"
	"database/sql"

	"github.com/Aifyel/petCryptoCurrency/internal/fetcher/domain"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}

func (r *CurrencyRepository) GetCurrencies(
	ctx context.Context,
) ([]domain.Currency, error) {
	return nil, nil
}

func (r *CurrencyRepository) SaveNewCurrency(
	ctx context.Context,
	currency domain.Currency,
) error {
	return nil
}

func (r *CurrencyRepository) Exists(
	ctx context.Context,
	currency domain.Currency,
) (bool, error) {
	return false, nil
}
