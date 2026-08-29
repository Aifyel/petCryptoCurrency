package usecase

import (
	"context"
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf("updater consumer: %w", entities.ErrInvalidParams)
	}

	if repo == nil {
		return nil, fmt.Errorf("updater repository: %w", entities.ErrInvalidParams)
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
				msgErrCh <- fmt.Errorf("updater ReadMessage: %w, %w", entities.ErrMessagingFailure, err)
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
				logErr := fmt.Errorf("update rate: %w, %w", entities.ErrInvalidMessage, err)
				log.Printf("%v\n", logErr)
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
		return fmt.Errorf("updater Save: %w, %w", entities.ErrRepositoryFailure, err)
	}

	err = u.consumer.CommitMessages(ctx, *msgs...)
	if err != nil {
		return fmt.Errorf("updater CommitMessages: %w, %w", entities.ErrMessagingFailure, err)
	}

	*rates = (*rates)[:0]
	*msgs = (*msgs)[:0]

	return nil
}
