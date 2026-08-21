package stub

import (
	"context"
	"sync"
)

type Repository struct {
	mu         sync.Mutex
	currencies []string
}

func NewRepository(initial []string) *Repository {
	c := make([]string, len(initial))
	copy(c, initial)
	return &Repository{currencies: c}
}

func (r *Repository) GetCurrencies(ctx context.Context) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]string, len(r.currencies))
	copy(out, r.currencies)
	return out, nil
}

func (r *Repository) SaveNewCurrency(ctx context.Context, currency []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.currencies = append(r.currencies, currency)
	return nil
}
