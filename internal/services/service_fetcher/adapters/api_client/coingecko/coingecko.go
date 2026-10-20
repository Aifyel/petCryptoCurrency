package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
)

const (
	vsCurrency = "usd"
)

type Client struct {
	httpClient  *http.Client
	baseURL     string
	minInterval time.Duration
	mu          sync.Mutex
	lastRequest time.Time
	logger      *zap.Logger
}

func NewClient(url string, minInterval time.Duration, logger *zap.Logger) (*Client, error) {
	if strings.TrimSpace(url) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko client: invalid url")
	}
	if minInterval <= 0 {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko client: invalid min interval")
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL:     url,
		minInterval: minInterval,
		lastRequest: time.Time{},
		logger:      logger,
	}, nil
}

func (c *Client) Fetch(ctx context.Context, currencies []string) ([]entities.CurrencyRate, error) {
	c.mu.Lock()
	elapsed := time.Since(c.lastRequest)
	if elapsed < c.minInterval {
		time.Sleep(c.minInterval - elapsed)
	}
	c.lastRequest = time.Now()
	c.mu.Unlock()

	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko: parse base url")
	}

	query := parsedURL.Query()
	query.Set("ids", strings.Join(currencies, ","))
	query.Set("vs_currencies", vsCurrency)
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko: build request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(entities.ErrClientFailure, "coingecko: http request")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(entities.ErrClientFailure, fmt.Sprintf("coingecko: unexpected status %d", resp.StatusCode))
	}

	var priceResp map[string]map[string]float64
	err = json.NewDecoder(resp.Body).Decode(&priceResp)
	if err != nil {
		return nil, errors.Wrap(entities.ErrClientFailure, "coingecko: decode response")
	}

	rates := make([]entities.CurrencyRate, 0)
	for id, prices := range priceResp {
		price, ok := prices[vsCurrency]
		if !ok {
			continue
		}

		rate, err := entities.NewCurrencyRate(id, price, time.Now().UTC())
		if err != nil {
			c.logger.Error("coingecko: new currency rate failed",
				zap.String("currency", id),
				zap.Float64("price", price),
				zap.Error(err),
			)
			continue
		}

		rates = append(rates, *rate)
	}

	return rates, nil
}
