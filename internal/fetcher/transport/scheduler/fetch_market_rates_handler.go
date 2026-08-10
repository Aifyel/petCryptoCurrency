package scheduler

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/fetcher/usecase"
)

type FetchMarketRatesHandler struct {
	Handler usecase.FetchMarketRatesUseCase
}

func NewFetchMarketRatesHandler(Handler usecase.FetchMarketRatesUseCase) *FetchMarketRatesHandler {
	return &FetchMarketRatesHandler{
		Handler: Handler,
	}
}

func (h *FetchMarketRatesHandler) Start(ctx context.Context) {
}
