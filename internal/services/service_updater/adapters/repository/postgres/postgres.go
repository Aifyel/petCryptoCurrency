package postgres

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, connString string) (*Repository, error) {
	if connString == "" {
	}
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (r *Repository) Save(ctx context.Context, currency []entities.CurrencyRate) error {
	builder := sq.Insert("rates").Columns("currency", "price", "fetched_at")
	for _, c := range currency {
		builder = builder.Values(c.Currency, c.Price, c.FetchedAt)
	}
	sql, args, err := builder.Suffix("ON CONFLICT (currency, fetched_at) DO NOTHING").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(entities.ErrInvalidParams, "repository: builder insert rates")
	}
	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(entities.ErrRepositoryFailure, "repository: save rates")
	}
	return nil
}
