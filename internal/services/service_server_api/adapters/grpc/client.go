package grpc

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petcrypto-proto/gen"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Client struct {
	client gen.FetcherServiceClient
}

func NewClient(conn *grpc.ClientConn) (*Client, error) {
	if conn == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gRPC client: connection is nil")
	}

	return &Client{
		client: gen.NewFetcherServiceClient(conn),
	}, nil
}

func (c *Client) Fetch(ctx context.Context, currency []string) ([]entities.CurrencyRate, error) {
	resp, err := c.client.FetchNewRates(ctx, &gen.FetchNewRatesRequest{Currencies: currency})
	if err != nil {
		return nil, errors.Wrap(entities.ErrClientFailure, "gRPC client: fetching currencies")
	}

	rates := make([]entities.CurrencyRate, 0)
	for _, rate := range resp.GetRates() {
		rates = append(rates, entities.CurrencyRate{
			Currency:  rate.GetCurrency(),
			Price:     rate.GetPrice(),
			FetchedAt: rate.GetFetchedAt().AsTime(),
		})
	}

	return rates, nil
}
