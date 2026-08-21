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

func NewCurrencyRate(currency string, price float64, fetchedAt time.Time) (*CurrencyRate, error) {
	if currency == "" {
		return nil, fmt.Errorf("invalid currency")
	} //...
	return &CurrencyRate{
		Currency: currency,
	}, nil
}
