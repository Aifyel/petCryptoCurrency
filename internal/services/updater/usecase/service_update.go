package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/segmentio/kafka-go"
)

const (
	maxSize      = 50
	flushTimeout = 5 * time.Second
)

type UpdateService struct {
	repo     RateRepository
	consumer Consumer
}

func NewUpdateService(repo RateRepository, consumer Consumer) (*UpdateService, error) {
	if consumer == nil {
		return nil, errors.New("consumer cannot be nil")
	}

	if repo == nil {
		return nil, errors.New("repo cannot be nil")
	}

	return &UpdateService{
		repo:     repo,
		consumer: consumer,
	}, nil
}

func (u *UpdateService) Execute(ctx context.Context) error {
	msgCh := make(chan kafka.Message)
	msgErrCh := make(chan error)

	go func() {
		for {
			msg, err := u.consumer.ReadMessage(ctx)
			if err != nil {
				msgErrCh <- err
				return
			}

			select {
			case msgCh <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	ticker := time.NewTicker(flushTimeout)
	defer ticker.Stop()

	rates := make([]entities.CurrencyRate, 0, maxSize)
	msgs := make([]kafka.Message, 0, maxSize)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-msgErrCh:
			return err
		case msg := <-msgCh:
			var rate entities.CurrencyRate
			err := json.Unmarshal(msg.Value, &rate)
			if err != nil {
				log.Printf("failed to unmarshal rate json: %v", err)
				continue
			}

			rates = append(rates, rate)
			msgs = append(msgs, msg)

			if len(msgs) >= maxSize {
				ctxWithTimeoutFlush, cancelFlush := context.WithTimeout(ctx, 1*time.Second)

				err = u.flush(ctxWithTimeoutFlush, &rates, &msgs)
				cancelFlush()
				if err != nil {
					log.Printf("failed to flush batch: %v", err)
					continue
				}
			}

		case <-ticker.C:
			if len(msgs) > 0 {
				ctxWithTimeoutFlush, cancelFlush := context.WithTimeout(ctx, 1*time.Second)

				err := u.flush(ctxWithTimeoutFlush, &rates, &msgs)
				cancelFlush()
				if err != nil {
					log.Printf("failed to flush batch: %v", err)
					continue
				}
			}
		}
	}
}

func (u *UpdateService) flush(ctx context.Context, rates *[]entities.CurrencyRate, msgs *[]kafka.Message) error {
	err := u.repo.Save(ctx, *rates)
	if err != nil {
		return err
	}

	err = u.consumer.CommitMessages(ctx, *msgs...)
	if err != nil {
		return err
	}

	*rates = (*rates)[:0]
	*msgs = (*msgs)[:0]

	return nil
}
