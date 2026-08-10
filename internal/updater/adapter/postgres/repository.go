package postgres

import (
	"context"
	"database/sql"

	"github.com/Aifyel/petCryptoCurrency/internal/updater/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(
	ctx context.Context,
	snapshot domain.PriceSnapshot,
) error {
	return nil
}
