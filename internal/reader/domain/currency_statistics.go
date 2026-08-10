package domain

type Currency string

type CurrencyStatistics struct {
	Currency          Currency
	CurrentPrice      float64
	LowestPriceDaily  float64
	HighestPriceDaily float64
	ChangeHourPercent float64
}
