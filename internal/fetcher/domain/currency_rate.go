package domain

import (
	"time"
)

type Currency string

type CurrencyRate struct {
	Currency  Currency
	Price     float64
	FetchedAt time.Time
}

func (r CurrencyRate) Validate() error {
	if r.Currency == "" {
		return nil
	}

	if r.Price < 0 || r.Price == 0 {
		return nil
	}

	if r.FetchedAt.IsZero() {
		return nil
	}

	return nil
}
