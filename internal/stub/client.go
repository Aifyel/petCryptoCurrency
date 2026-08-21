package stub

import (
	"context"
	"time"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type Client struct {
	catalog map[string]float64
}

func NewClient() *Client {
	return &Client{
		catalog: map[string]float64{
			"USD":  90.5,
			"EUR":  98.3,
			"GBP":  115.2,
			"TRY":  2.7,
			"MEOW": 100,
		},
	}
}

func (c *Client) Fetch(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error) {
	rates := make([]entities.CurrencyRate, 0, len(currencies))

	for _, cur := range currencies {
		price, ok := c.catalog[cur]
		if !ok {
			continue
		}
		rates = append(rates, entities.CurrencyRate{
			Currency:  cur,
			Price:     price,
			FetchedAt: time.Now(),
		})
	}

	return rates, nil
}
