package entities

import (
	"fmt"
	"time"
)

type CurrencyRate struct {
	Currency  string
	Price     float64
	FetchedAt time.Time
}

func NewCurrencyRate(c string, p float64, f time.Time) (*CurrencyRate, error) {
	if c == "" {
		return nil, fmt.Errorf("currencyRate currency: %w", ErrInvalidParams)
	}
	if p < 0 {
		return nil, fmt.Errorf("currencyRate price: %w", ErrInvalidParams)
	}

	return &CurrencyRate{
		Currency:  c,
		Price:     p,
		FetchedAt: f,
	}, nil
}
