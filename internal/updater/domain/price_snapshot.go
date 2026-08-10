package domain

import (
	"time"
)

type Currency string

type PriceSnapshot struct {
	Currency   Currency
	Price      float64
	RecordedAt time.Time
}
