package app

import (
	"context"
	"log"

	"github.com/Aifyel/petCryptoCurrency/deploy/config"
	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_updater/adapters/repository/postgres"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_updater/ports/broker_consumer/kafka"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_updater/usecase"
	"github.com/pkg/errors"
)

func Run(ctx context.Context, cfg *config.Config) error {
	repo, err := postgres.NewRepository(ctx, cfg.DatabaseURL)
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
			log.Printf("updater_app: failed to close consumer: %v", err)
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
