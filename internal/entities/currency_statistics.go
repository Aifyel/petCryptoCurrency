package entities

import (
	"fmt"
)

type CurrencyStatistics struct {
	Currency          string
	CurrentPrice      float64
	LowestPriceDaily  float64
	HighestPriceDaily float64
	ChangeHourPercent float64
}

func NewCurrencyStatistics(c string, cp float64, lp float64, hp float64, chp float64) (*CurrencyStatistics, error) {
	if c == "" {
		return nil, fmt.Errorf("CurrencyStatistics currency: %w", ErrInvalidParams)
	}
	if cp < 0 {
		return nil, fmt.Errorf("CurrencyStatistics price: %w", ErrInvalidParams)
	}
	if lp < 0 {
		return nil, fmt.Errorf("CurrencyStatistics lowest price: %w", ErrInvalidParams)
	}
	if hp < 0 {
		return nil, fmt.Errorf("CurrencyStatistics highest price: %w", ErrInvalidParams)
	}
	if lp > hp {
		return nil, fmt.Errorf("CurrencyStatistics highest price: %w", ErrInvalidParams)
	}

	return &CurrencyStatistics{
		Currency:          c,
		CurrentPrice:      cp,
		LowestPriceDaily:  lp,
		HighestPriceDaily: hp,
		ChangeHourPercent: chp,
	}, nil
}
