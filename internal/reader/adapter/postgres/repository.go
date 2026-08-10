package postgres

import (
	"context"
	"database/sql"

	"github.com/Aifyel/petCryptoCurrency/internal/reader/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStatistics(
	ctx context.Context,
	currency domain.Currency,
) (domain.CurrencyStatistics, error) {
	return domain.CurrencyStatistics{}, nil
}
