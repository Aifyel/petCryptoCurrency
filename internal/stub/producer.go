package stub

import (
	"context"
	"log"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

type Producer struct{}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Produce(ctx context.Context, rates []entities.CurrencyRate) error {
	log.Printf("publish: %+v", rates)
	return nil
}
