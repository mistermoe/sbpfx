package sbpfx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mistermoe/httpr"
)

const (
	BaseURL      = "https://www.sbp.org.pk/assets/document"
	HTTPStatusOK = 200
	YearModulo   = 100 // For getting last 2 digits of year
	ratePrefix   = "/mark-to-market-revaluation-exchange-rate"
	pdfSignature = "%PDF" // PDF files begin with this magic header
)

// looksLikePDF reports whether content begins with the PDF file signature.
// The new /assets/document host serves a 200 with a non-PDF body for missing
// sheets instead of a 404, so the status code alone is not enough to tell a
// real rate sheet from a "not found" response.
func looksLikePDF(content []byte) bool {
	return bytes.HasPrefix(content, []byte(pdfSignature))
}

// newFormatStart is the first date whose sheet uses the current filename
// format (prefix + DD-month-YYYY, e.g. 03-july-2026). Earlier sheets use the
// legacy prefix + DD-Mon-YY style (e.g. 27-Aug-25).
var newFormatStart = time.Date(2026, time.July, 3, 0, 0, 0, 0, time.UTC)

// migrationOverrides holds the handful of transition-window sheets that SBP
// uploaded with irregular names while migrating to the new host. Keyed by
// YYYY-MM-DD. These follow neither the legacy nor the current convention, so
// they can only be resolved by lookup.
var migrationOverrides = map[string]string{
	"2026-06-29": "/29-Jun-26.pdf",              // no prefix
	"2026-06-30": "/30-Jun-26_1.pdf",            // no prefix, _1 suffix
	"2026-07-02": ratePrefix + "-02-Jul-26.pdf", // prefix but legacy date style
}

// ratePath builds the URL path for the exchange rate PDF of the given date.
//   - transition-window dates: looked up in migrationOverrides
//   - before newFormatStart: legacy prefix + DD-Mon-YY (e.g. 27-Aug-25)
//   - on/after newFormatStart: current prefix + DD-month-YYYY (e.g. 07-july-2026)
func ratePath(date time.Time) string {
	if override, ok := migrationOverrides[date.Format("2006-01-02")]; ok {
		return override
	}

	if date.Before(newFormatStart) {
		return fmt.Sprintf("%s-%02d-%s-%02d.pdf",
			ratePrefix,
			date.Day(),
			date.Format("Jan"),
			date.Year()%YearModulo, // Last 2 digits of year
		)
	}

	return fmt.Sprintf("%s-%02d-%s-%d.pdf",
		ratePrefix,
		date.Day(),
		strings.ToLower(date.Format("January")),
		date.Year(),
	)
}

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
	path := ratePath(date)

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

	if !looksLikePDF(content) {
		return nil, fmt.Errorf("PDF not found: no rate sheet available for path: %s", path)
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

	return fmt.Sprintf("%s%s", BaseURL, ratePath(date))
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
	urlPath := ratePath(date)

	resp, err := c.httpClient.Get(ctx, urlPath)
	if err != nil {
		return fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != HTTPStatusOK {
		return fmt.Errorf("PDF not found: status %d for path: %s", resp.StatusCode, urlPath)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read PDF: %w", err)
	}

	if !looksLikePDF(content) {
		return fmt.Errorf("PDF not found: no rate sheet available for path: %s", urlPath)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return fmt.Errorf("failed to write PDF to file %s: %w", path, err)
	}

	return nil
}
