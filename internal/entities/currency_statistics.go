package entities

type CurrencyStatistics struct {
	Currency          string
	CurrentPrice      float64
	LowestPriceDaily  float64
	HighestPriceDaily float64
	ChangeHourPercent float64
}
