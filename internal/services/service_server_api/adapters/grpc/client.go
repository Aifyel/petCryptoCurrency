package grpc

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/deploy/proto/gen"
	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Client struct {
	client gen.FetcherServiceClient
	logger *zap.Logger
}

func NewClient(conn *grpc.ClientConn, logger *zap.Logger) (*Client, error) {
	if conn == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gRPC client: connection is nil")
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	return &Client{
		client: gen.NewFetcherServiceClient(conn),
		logger: logger,
	}, nil
}

func (c *Client) Fetch(ctx context.Context, currency []string) ([]entities.CurrencyRate, error) {
	resp, err := c.client.FetchNewRates(ctx, &gen.FetchNewRatesRequest{Currencies: currency})
	if err != nil {
		c.logger.Error("gRPC client: failed to fetch currencies",
			zap.Error(err),
			zap.Strings("currencies", currency),
		)
		return nil, errors.Wrap(entities.ErrClientFailure, "gRPC client: fetching currencies")
	}

	rates := make([]entities.CurrencyRate, 0)
	for _, rate := range resp.GetRates() {
		rates = append(rates, entities.CurrencyRate{
			Currency:  rate.GetCurrency(),
			Price:     rate.GetPrice(),
			FetchedAt: rate.GetFetchedAt().AsTime().UTC(),
		})
	}

	return rates, nil
}
