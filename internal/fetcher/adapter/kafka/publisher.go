package kafka

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/fetcher/domain"
	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	broker *kafka.Writer
}

func NewPublisher(broker *kafka.Writer) *Publisher {
	return &Publisher{
		broker: broker,
	}
}

func (p *Publisher) Publish(
	ctx context.Context,
	rate domain.CurrencyRate,
) error {
	return nil
}
