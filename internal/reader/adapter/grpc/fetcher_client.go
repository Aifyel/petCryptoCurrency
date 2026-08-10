package grpc

import (
	"context"

	fetcher "github.com/Aifyel/petCryptoCurrency/internal/fetcher/domain"
	"github.com/Aifyel/petCryptoCurrency/internal/reader/domain"
)

type FetcherClient struct {
	client FetcherServiceClient //proto gen
}

func NewFetcherClient(client FetcherServiceClient) *FetcherClient {
	return &FetcherClient{client: client}
}

func (c *FetcherClient) FetchCurrency(ctx context.Context, currency domain.Currency) (fetcher.CurrencyRate, error) {
	return fetcher.CurrencyRate{}, nil
}
