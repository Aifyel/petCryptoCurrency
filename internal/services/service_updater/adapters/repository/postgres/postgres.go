package postgres

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Repository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewRepository(ctx context.Context, connString string, logger *zap.Logger) (*Repository, error) {
	if connString == "" {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko: empty connection string")
	}
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, err.Error())
	}

	return &Repository{
		pool:   pool,
		logger: logger,
	}, nil
}

func (r *Repository) Save(ctx context.Context, rates []entities.CurrencyRate) error {
	if len(rates) == 0 {
		return nil
	}

	builder := sq.Insert("rates").Columns("currency", "price", "fetched_at")

	for _, c := range rates {
		builder = builder.Values(c.Currency, c.Price, c.FetchedAt)
	}

	sql, args, err := builder.Suffix("ON CONFLICT (currency, fetched_at) DO NOTHING").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(entities.ErrInvalidParams, "repository: builder insert rates")
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error("repository: failed to save rates",
			zap.Error(err),
			zap.Int("rates_count", len(rates)),
		)
		return errors.Wrap(entities.ErrRepositoryFailure, "repository: save rates")
	}

	r.logger.Debug("rates successfully saved to database",
		zap.Int("rates_count", len(rates)),
	)

	return nil
}

func (r *Repository) Close() {
	r.pool.Close()
}
