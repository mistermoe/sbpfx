package sbpfx

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mistermoe/httpr"
)

const (
	BaseURL      = "https://www.sbp.org.pk/ecodata/rates/m2m"
	HTTPStatusOK = 200
	YearModulo   = 100 // For getting last 2 digits of year
)

type Client struct {
	httpClient *httpr.Client
}

func New(options ...httpr.ClientOption) *Client {
	opts := append(
		[]httpr.ClientOption{
			httpr.BaseURL(BaseURL),
		},
		options...,
	)

	return &Client{
		httpClient: httpr.NewClient(opts...),
	}
}

func (c *Client) GetExchangeRates(ctx context.Context, opts ...Option) (map[Currency]*ExchangeRate, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	date := cfg.date
	path := fmt.Sprintf("/%d/%s/%02d-%s-%02d.pdf",
		date.Year(),
		date.Format("Jan"),
		date.Day(),
		date.Format("Jan"),
		date.Year()%YearModulo) // Last 2 digits of year

	fullURL := fmt.Sprintf("%s%s", BaseURL, path)

	resp, err := c.httpClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != HTTPStatusOK {
		return nil, fmt.Errorf("PDF not found: status %d for path: %s", resp.StatusCode, path)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	return parsePDFContent(content, date, fullURL)
}

func (c *Client) GetExchangeRate(ctx context.Context, currency Currency, opts ...Option) (*ExchangeRate, error) {
	rates, err := c.GetExchangeRates(ctx, opts...)
	if err != nil {
		return nil, err
	}

	rate, exists := rates[currency]
	if !exists {
		return nil, fmt.Errorf("exchange rate for %s not found", currency)
	}

	return rate, nil
}

func (c *Client) GetUrl(opts ...Option) string {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			// For this method, we'll ignore errors and use default config
			// since it returns a string and can't propagate the error
			continue
		}
	}

	date := cfg.date

	return fmt.Sprintf("%s/%d/%s/%02d-%s-%02d.pdf",
		BaseURL,
		date.Year(),
		date.Format("Jan"),
		date.Day(),
		date.Format("Jan"),
		date.Year()%YearModulo, // Last 2 digits of year
	)
}

// DownloadRateSheet downloads the exchange rate PDF to the specified file path.
func (c *Client) DownloadRateSheet(ctx context.Context, path string, opts ...Option) error {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return fmt.Errorf("failed to apply option: %w", err)
		}
	}

	date := cfg.date
	urlPath := fmt.Sprintf("/%d/%s/%02d-%s-%02d.pdf",
		date.Year(),
		date.Format("Jan"),
		date.Day(),
		date.Format("Jan"),
		date.Year()%YearModulo) // Last 2 digits of year

	resp, err := c.httpClient.Get(ctx, urlPath)
	if err != nil {
		return fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != HTTPStatusOK {
		return fmt.Errorf("PDF not found: status %d for path: %s", resp.StatusCode, urlPath)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write PDF to file %s: %w", path, err)
	}

	return nil
}
