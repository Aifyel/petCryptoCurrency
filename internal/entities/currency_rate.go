package entities

import (
	"time"
)

type CurrencyRate struct {
	Currency  string
	Price     float64
	FetchedAt time.Time
}
