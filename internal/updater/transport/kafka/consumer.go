package kafka

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/updater/domain"
	"github.com/segmentio/kafka-go"
)

type RatePersister interface {
	Execute(
		ctx context.Context,
		snapshot domain.PriceSnapshot,
	) error
}

type Consumer struct {
	reader    *kafka.Reader
	persister RatePersister
}

func NewConsumer(
	reader *kafka.Reader,
	persister RatePersister,
) *Consumer {
	return &Consumer{
		reader:    reader,
		persister: persister,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	return nil
}
