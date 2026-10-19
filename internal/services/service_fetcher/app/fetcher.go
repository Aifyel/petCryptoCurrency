package app

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	"github.com/Aifyel/petCryptoCurrency/deploy/proto/gen"
	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/adapters/api_client/coingecko"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/adapters/broker_producer/kafka"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/adapters/repository/postgres"
	grpcport "github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/ports/grpc"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/usecase"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg *config.Config) error {
	repo, err := postgres.NewRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		return errors.Wrap(entities.ErrRepositoryFailure, "fetcher_app: failed to connect to postgres")
	}
	defer repo.Close()

	client, err := coingecko.NewClient(cfg.CoingeckoURL, cfg.CoingeckoInterval)
	if err != nil {
		return errors.Wrap(entities.ErrClientFailure, "fetcher_app: failed to connect to coingecko")
	}

	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		return errors.Wrap(err, "fetcher_app: failed to create producer")
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("fetcher_app: failed to close producer: %v", err)
		}
	}()

	fetchService, err := usecase.NewFetchService(repo, client, producer)
	if err != nil {
		return errors.Wrap(entities.ErrClientFailure, "fetcher_app: failed to create fetch service")
	}

	go func() {
		initCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := fetchService.UpdateRates(initCtx); err != nil {
			log.Printf("fetcher_app: initial update rates failed: %v", err)
		}
	}()

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

	ctxWC, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunCron(ctxWC, fetchService)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctxWC.Done()
		grpcServer.GracefulStop()
	}()

	serveErr := grpcServer.Serve(listener)

	cancel()
	wg.Wait()

	if serveErr != nil {
		return errors.Wrap(serveErr, "fetcher_app: grpc server failed")
	}

	return nil
}

func RunCron(ctx context.Context, service *usecase.FetchService) {
	crn := cron.New()

	updateFunc := func() {
		uCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := service.UpdateRates(uCtx)
		if err != nil {
			log.Printf("cron_fetcher: failed to update rates: %v", err)
		}
	}

	_, err := crn.AddFunc("@every 5m", updateFunc)
	if err != nil {
		log.Printf("cron_fetcher: failed to add func: %v", err)
		return
	}

	crn.Start()

	<-ctx.Done()

	stopCtx := crn.Stop()

	<-stopCtx.Done()
}
