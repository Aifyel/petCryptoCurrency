package entities

import (
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(ErrInvalidParams, "CurrencyStatistics currency:")
	}
	if cp < 0 {
		return nil, errors.Wrap(ErrInvalidParams, "CurrencyStatistics currentPrice:")
	}
	if lp < 0 {
		return nil, errors.Wrap(ErrInvalidParams, "CurrencyStatistics lowestPrice:")
	}
	if hp < 0 {
		return nil, errors.Wrap(ErrInvalidParams, "CurrencyStatistics highestPrice:")
	}
	if lp > hp {
		return nil, errors.Wrap(ErrInvalidParams, "lowerPrice should not be greater than highPrice:")
	}

	return &CurrencyStatistics{
		Currency:          c,
		CurrentPrice:      cp,
		LowestPriceDaily:  lp,
		HighestPriceDaily: hp,
		ChangeHourPercent: chp,
	}, nil
}
