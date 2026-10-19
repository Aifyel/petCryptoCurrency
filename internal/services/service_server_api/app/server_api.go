package app

import (
	"context"
	"time"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	servicegrpc "github.com/Aifyel/petCryptoCurrency/internal/services/service_server_api/adapters/grpc"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_server_api/adapters/repository/postgres"
	ginservice "github.com/Aifyel/petCryptoCurrency/internal/services/service_server_api/ports/http/public/gin"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_server_api/usecase"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(ctx context.Context, cfg *config.Config) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create logger")
	}
	defer logger.Sync()

	repo, err := postgres.NewRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to connect to postgres")
	}
	defer repo.Close()

	conn, err := grpc.NewClient(cfg.GRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create grpc client")
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logger.Error("server_api: unable to close grpc client connection", zap.Error(err))
		}
	}(conn)

	client, err := servicegrpc.NewClient(conn)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create grpc client")
	}

	service, err := usecase.NewGetCurrencyStatistics(repo, client)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create new service")
	}

	server, err := ginservice.NewServer(cfg.HTTPPort, 5*time.Second, service, logger)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create gin server")
	}

	err = server.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to start gin server")
	}

	return nil
}
