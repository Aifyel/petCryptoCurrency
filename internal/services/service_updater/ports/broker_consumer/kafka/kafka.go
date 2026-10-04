package kafka

import (
	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/Aifyel/petCryptoCurrency/internal/services/service_fetcher/usecase"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func NewConsumer(brokers []string, groupID string) (*kafka.Reader, error) {
	if len(brokers) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParams, "consumer: empty list of brokers")
	}
	if groupID == "" {
		return nil, errors.Wrap(entities.ErrInvalidParams, "consumer: empty groupID")
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: []string{string(usecase.TopicScheduled), string(usecase.TopicNew)},
	}), nil
}
