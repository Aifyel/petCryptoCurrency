package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

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
}

func NewClient(url string, minInterval time.Duration) (*Client, error) {
	if strings.TrimSpace(url) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko client: invalid url")
	}
	if minInterval <= 0 {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko client: invalid min interval")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		baseURL:     url,
		minInterval: minInterval,
		lastRequest: time.Time{},
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

	ids := strings.Join(currencies, ",")
	url := fmt.Sprintf(c.baseURL+"?ids=%s&vs_currencies=%s", ids, vsCurrency) //Изменить на net/url пакет

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "coingecko: build request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(entities.ErrClientFailure, "coingecko: http request")
	}

	defer resp.Body.Close()

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

		rate, err := entities.NewCurrencyRate(id, price, time.Now())
		if err != nil {
			log.Printf("%v", errors.Wrap(entities.ErrInvalidParams, "coingecko: new currency rate"))
			continue
		}

		rates = append(rates, *rate)
	}

	return rates, nil
}
