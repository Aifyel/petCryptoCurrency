package grpc

import "github.com/Aifyel/petCryptoCurrency/internal/fetcher/usecase"

type FetchCurrencyHandler struct {
	Handler usecase.FetchCurrencyGRPC
}

func NewFetchCurrencyHandler(Handler usecase.FetchCurrencyGRPC) *FetchCurrencyHandler {
	return &FetchCurrencyHandler{
		Handler: Handler,
	}
}
