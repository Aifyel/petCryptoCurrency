package grpcport

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/deploy/proto/gen"
	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Executor interface {
	FetchNewRates(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error)
}
type Server struct {
	executor Executor
	logger   *zap.Logger
	gen.UnimplementedFetcherServiceServer
}

func NewServer(executor Executor, logger *zap.Logger) (*Server, error) {
	if executor == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gRPC server: nil executor")
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	return &Server{
		executor: executor,
		logger:   logger,
	}, nil
}

func (s *Server) FetchNewRates(ctx context.Context, req *gen.FetchNewRatesRequest) (*gen.FetchNewRatesResponse, error) {
	if len(req.GetCurrencies()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "currencies list is empty")
	}

	currencies, err := s.executor.FetchNewRates(ctx, req.GetCurrencies())
	if err != nil {
		s.logger.Error("gRPC server: failed to fetch currencies",
			zap.Error(err),
			zap.Strings("currencies", req.GetCurrencies()),
		)
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
