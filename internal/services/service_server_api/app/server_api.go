package app

import (
	"context"
	"time"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	"github.com/Aifyel/petCryptoCurrency/internal/logger"
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
	if err := logger.InitLogger(false); err != nil {
		return errors.Wrap(err, "failed to init logger")
	}
	defer logger.Sync()

	logger.Log.Info("service_api is starting up...")

	repo, err := postgres.NewRepository(ctx, cfg.DatabaseURL, logger.Log)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to connect to postgres")
	}
	defer repo.Close()

	conn, err := grpc.NewClient(cfg.GRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create grpc client")
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			logger.Log.Error("server_api: unable to close grpc client connection", zap.Error(err))
		}
	}(conn)

	client, err := servicegrpc.NewClient(conn, logger.Log)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create grpc client")
	}

	service, err := usecase.NewGetCurrencyStatistics(repo, client)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create new service")
	}

	server, err := ginservice.NewServer(cfg.HTTPPort, 5*time.Second, service, logger.Log)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to create gin server")
	}

	err = server.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "server_api: unable to start gin server")
	}

	return nil
}
