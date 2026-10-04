package app

import (
	"context"
	"log"
	"net"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/adapters/api_client/coingecko"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/adapters/broker_producer/kafka"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/adapters/repository/postgres"
	grpcport "github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/ports/grpc"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/usecase"
	"github.com/Aifyel/petcrypto-proto/gen"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg *config.Config) error {
	repo, err := postgres.NewRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		return errors.Wrap(entities.ErrRepositoryFailure, "fetcher_app: failed to connect to postgres")
	}

	client, err := coingecko.NewClient(cfg.CoingeckoURL, cfg.CoingeckoInterval)
	if err != nil {
		return errors.Wrap(entities.ErrClientFailure, "fetcher_app: failed to connect to coingecko")
	}

	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		return errors.Wrap(err, "fetcher_app: failed to create producer")
	}
	defer producer.Close()

	fetchService, err := usecase.NewFetchService(repo, client, producer)
	if err != nil {
		return errors.Wrap(entities.ErrClientFailure, "fetcher_app: failed to create fetch service")
	}

	grpcHandler, err := grpcport.NewServer(fetchService)
	if err != nil {
		return errors.Wrap(err, "fetcher_app: failed to create grpc handler")
	}
	grpcServer := grpc.NewServer()
	gen.RegisterFetcherServiceServer(grpcServer, grpcHandler)

	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return errors.Wrap(err, "fetcher_app: failed to listen tcp port")
	}

	go RunCron(ctx, fetchService)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	err = grpcServer.Serve(listener)
	if err != nil {
		return errors.Wrap(err, "fetcher_app: grpc server failed")
	}

	return nil
}

func RunCron(ctx context.Context, service *usecase.FetchService) {
	crn := cron.New()

	updateFunc := func() {
		err := service.UpdateRates(ctx)
		if err != nil {
			log.Printf("cron_fetcher: failed to update rates: %v", err)
		}
	}

	_, err := crn.AddFunc("@every 30s", updateFunc)
	if err != nil {
		log.Printf("cron_fetcher: failed to add func: %v", err)
		return
	}

	go func() {
		<-ctx.Done()
		stopCtx := crn.Stop()
		<-stopCtx.Done()
	}()

	crn.Run()
}
