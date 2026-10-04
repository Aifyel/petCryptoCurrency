package postgres

import (
	"context"
	"log"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const statsReq = `
WITH
daily_stats AS (
	SELECT
		currency,
		MIN(price) AS lowest_price,
		MAX(price) AS highest_price
	FROM rates
	WHERE currency = ANY($1)
	  AND fetched_at >= NOW() - INTERVAL '24 hours'
	GROUP BY currency
),
current_price AS (
	SELECT DISTINCT ON (currency)
		currency,
		price AS current_price,
		fetched_at
	FROM rates
	WHERE currency = ANY($1)
	ORDER BY currency, fetched_at DESC
),
hour_ago_price AS (
	SELECT DISTINCT ON (currency)
		currency,
		price AS price_hour_ago
	FROM rates
	WHERE currency = ANY($1)
	  AND fetched_at <= NOW() - INTERVAL '1 hour'
	ORDER BY currency, fetched_at DESC
)
SELECT
	cp.currency,
	cp.current_price,
	COALESCE(ds.lowest_price, cp.current_price) AS lowest_price,
	COALESCE(ds.highest_price, cp.current_price) AS highest_price,
	COALESCE(
		(cp.current_price - hap.price_hour_ago) / NULLIF(hap.price_hour_ago, 0) * 100,
		0
	) AS change_hour_percent
FROM current_price cp
LEFT JOIN daily_stats ds ON cp.currency = ds.currency
LEFT JOIN hour_ago_price hap ON cp.currency = hap.currency
`

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

func (r *Repository) GetStatistics(ctx context.Context, currencies []string) ([]entities.CurrencyStatistics, error) {
	rows, err := r.pool.Query(ctx, statsReq, currencies)
	if err != nil {
		return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: query currencies")
	}

	defer rows.Close()

	var stats []entities.CurrencyStatistics
	for rows.Next() {
		var (
			currency                                     string
			currentPrice, lowest, highest, changePercent float64
		)
		err = rows.Scan(&currency, &currentPrice, &lowest, &highest, &changePercent)
		if err != nil {
			return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: rows scan")
		}

		currencyStatistics, err := entities.NewCurrencyStatistics(currency, currentPrice, lowest, highest, changePercent)
		if err != nil {
			log.Printf("%v", errors.Wrap(entities.ErrInvalidParams, "repository: new currency statistics"))
			continue
		}

		stats = append(stats, *currencyStatistics)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(entities.ErrRepositoryFailure, "repository: rows iteration")
	}

	return stats, nil
}
