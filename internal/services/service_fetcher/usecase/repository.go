package usecase

import "context"

type Repository interface {
	GetCurrencies(ctx context.Context) ([]string, error)
	SaveNewCurrency(ctx context.Context, currency []string) error
}
