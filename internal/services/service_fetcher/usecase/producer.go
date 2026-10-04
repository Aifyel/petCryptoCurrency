package usecase

import (
	"context"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type Topic string

const (
	TopicScheduled Topic = "rates.scheduled"
	TopicNew       Topic = "rates.new"
)

type Producer interface {
	Produce(ctx context.Context, rate []entities.CurrencyRate, topic Topic) error
}
