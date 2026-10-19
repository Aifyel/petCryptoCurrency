package kafka

import (
	"context"
	"encoding/json"
	"log"

	errorsbase "errors"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/usecase"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	scheduled *kafka.Writer
	new       *kafka.Writer
}

func NewProducer(brokers []string) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParams, "broker_producer: empty list of brokers")
	}

	return &Producer{
		scheduled: &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: string(usecase.TopicScheduled), Balancer: &kafka.LeastBytes{}},
		new:       &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: string(usecase.TopicNew), Balancer: &kafka.LeastBytes{}},
	}, nil
}

func (p *Producer) Produce(ctx context.Context, rate []entities.CurrencyRate, topic usecase.Topic) error {
	var writer *kafka.Writer
	switch topic {
	case usecase.TopicNew:
		writer = p.new
	case usecase.TopicScheduled:
		writer = p.scheduled
	default:
		return errors.Wrap(entities.ErrInvalidParams, "broker_producer: invalid topic")
	}

	msgs := make([]kafka.Message, 0, len(rate))
	for _, r := range rate {
		data, err := json.Marshal(r)
		if err != nil {
			log.Printf("%v: failed to marshal rate", entities.ErrMessagingFailure)
			continue
		}
		msgs = append(msgs, kafka.Message{Value: data})
	}

	err := writer.WriteMessages(ctx, msgs...)
	if err != nil {
		return errors.Wrap(entities.ErrMessagingFailure, "broker_producer: write messages")
	}

	return nil
}

func (p *Producer) Close() error {
	var errs error

	if err := p.scheduled.Close(); err != nil {
		errs = errorsbase.Join(errs, errors.Wrap(err, "broker_producer: failed to close writer: scheduled"))
	}

	if err := p.new.Close(); err != nil {
		errs = errorsbase.Join(errs, errors.Wrap(err, "broker_producer: failed to close writer: new"))
	}

	return errs
}
