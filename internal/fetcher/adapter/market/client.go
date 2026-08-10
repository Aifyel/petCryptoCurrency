package market

import (
	"context"
	"net/http"

	"github.com/Aifyel/petCryptoCurrency/internal/fetcher/domain"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Fetch(
	ctx context.Context,
	currencies []domain.Currency,
) ([]domain.CurrencyRate, error) {
	return nil, nil
}
