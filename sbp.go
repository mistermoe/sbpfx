package sbpfx

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mistermoe/httpr"
)

const (
	BaseURL = "https://www.sbp.org.pk/ecodata/rates/m2m"
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
		opt(cfg)
	}

	date := cfg.date
	// Build URL path - relative to base URL
	path := fmt.Sprintf("/%d/%s/%02d-%s-%02d.pdf",
		date.Year(),
		date.Format("Jan"),
		date.Day(),
		date.Format("Jan"),
		date.Year()%100) // Last 2 digits of year

	// Build full URL for reference
	fullURL := fmt.Sprintf("%s%s", BaseURL, path)

	// Download PDF using httpr
	resp, err := c.httpClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("PDF not found: status %d for path: %s", resp.StatusCode, path)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	return parsePDFContent(content, date, fullURL)
}

// GetExchangeRate returns the exchange rate for a specific currency
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

// GetUrl returns the fully qualified PDF URL for the specified date
func (c *Client) GetUrl(opts ...Option) string {
	// Apply options to default config
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	date := cfg.date

	return fmt.Sprintf("%s/%d/%s/%02d-%s-%02d.pdf",
		BaseURL,
		date.Year(),
		date.Format("Jan"),
		date.Day(),
		date.Format("Jan"),
		date.Year()%100, // Last 2 digits of year
	)
}

// DownloadRateSheet downloads the exchange rate PDF to the specified file path
func (c *Client) DownloadRateSheet(ctx context.Context, path string, opts ...Option) error {
	// Apply options to default config
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	date := cfg.date
	// Build URL path - relative to base URL
	urlPath := fmt.Sprintf("/%d/%s/%02d-%s-%02d.pdf",
		date.Year(),
		date.Format("Jan"),
		date.Day(),
		date.Format("Jan"),
		date.Year()%100) // Last 2 digits of year

	// Download PDF using httpr
	resp, err := c.httpClient.Get(ctx, urlPath)
	if err != nil {
		return fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("PDF not found: status %d for path: %s", resp.StatusCode, urlPath)
	}

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write PDF to file %s: %w", path, err)
	}

	return nil
}
