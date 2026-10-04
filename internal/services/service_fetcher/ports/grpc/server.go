package grpc

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petcrypto-proto/gen"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Executor interface {
	FetchNewRates(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error)
}
type Server struct {
	gen.UnimplementedFetcherServiceServer
	executor Executor
}

func NewServer(executor Executor) (*Server, error) {
	if executor == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gRPC server: nil executor")
	}
	return &Server{
		executor: executor,
	}, nil
}

func (s *Server) FetchNewRates(ctx context.Context, req *gen.FetchNewRatesRequest) (*gen.FetchNewRatesResponse, error) {
	if len(req.GetCurrencies()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "currencies list is empty")
	}

	currencies, err := s.executor.FetchNewRates(ctx, req.GetCurrencies())
	if err != nil {
		return nil, status.Error(codes.Internal, "gRPC server: failed to fetch currencies")
	}

	var currencyRate []*gen.CurrencyRate

	for _, currency := range currencies {
		fetchTime := timestamppb.New(currency.FetchedAt)

		currencyRate = append(currencyRate, &gen.CurrencyRate{
			Currency:  currency.Currency,
			Price:     currency.Price,
			FetchedAt: fetchTime,
		})
	}

	resp := &gen.FetchNewRatesResponse{
		Rates: currencyRate,
	}

	return resp, nil
}
