package app

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/logger"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_updater/adapters/repository/postgres"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_updater/ports/broker_consumer/kafka"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_updater/usecase"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg *config.Config) error {
	if err := logger.InitLogger(false); err != nil {
		return errors.Wrap(err, "failed to init logger")
	}
	defer logger.Sync()

	logger.Log.Info("updater-service is starting up...")

	repo, err := postgres.NewRepository(ctx, cfg.DatabaseURL, logger.Log)
	if err != nil {
		return errors.Wrap(entities.ErrRepositoryFailure, "updater_app: failed to connect to postgres")
	}
	defer repo.Close()

	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID)
	if err != nil {
		return errors.Wrap(entities.ErrMessagingFailure, "updater_app: failed to create kafka consumer")
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Log.Error("updater_app: failed to close consumer", zap.Error(err))
		}
	}()

	updateService, err := usecase.NewUpdateService(repo, consumer)
	if err != nil {
		return errors.Wrap(err, "updater_app: failed to create update service")
	}

	err = updateService.Execute(ctx)
	if err != nil {
		return errors.Wrap(err, "updater_app: failed to run Execute")
	}

	return nil
}
