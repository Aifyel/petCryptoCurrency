package entities

import (
	"time"

	"github.com/pkg/errors"
)

type CurrencyRate struct {
	Currency  string
	Price     float64
	FetchedAt time.Time
}

func NewCurrencyRate(c string, p float64, t time.Time) (*CurrencyRate, error) {
	if c == "" {
		return nil, errors.Wrap(ErrInvalidParams, "currencyRate currency:")
	}
	if p < 0 {
		return nil, errors.Wrap(ErrInvalidParams, "currencyRate price:")
	}

	return &CurrencyRate{
		Currency:  c,
		Price:     p,
		FetchedAt: t.UTC(),
	}, nil
}
