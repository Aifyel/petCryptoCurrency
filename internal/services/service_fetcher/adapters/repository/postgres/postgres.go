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
		return nil, errors.Wrap(entities.ErrInvalidParams, "repository: invalid connection string")
	}
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: failure pool create")
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: failure pool ping")
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (r *Repository) GetCurrencies(ctx context.Context) ([]string, error) {
	sql, args, err := sq.Select("currency").
		From("currencies").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "repository: GetCurrencies")
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: query currencies")
	}
	defer rows.Close()

	var currencies []string
	for rows.Next() {
		var currency string
		err := rows.Scan(&currency)
		if err != nil {
			return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: scan currency")
		}
		currencies = append(currencies, currency)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: rows iteration")
	}

	return currencies, nil
}

func (r *Repository) SaveNewCurrency(ctx context.Context, currencies []string) error {
	builder := sq.Insert("currencies").Columns("currency")
	for _, c := range currencies {
		builder = builder.Values(c)
	}
	sql, args, err := builder.Suffix("ON CONFLICT (currency) DO NOTHING").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(entities.ErrInvalidParams, "repository: builder insert currencies")
	}
	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(entities.ErrRepositoryFailure, "repository: save currencies")
	}

	return nil
}
