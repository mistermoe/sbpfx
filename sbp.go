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
	BaseURL        = "https://www.sbp.org.pk/assets/document"
	HTTPStatusOK   = 200
	YearModulo     = 100 // For getting last 2 digits of year
	ratePrefix     = "/mark-to-market-revaluation-exchange-rate"
	pdfSignature   = "%PDF"            // PDF files begin with this magic header
	pdfContentType = "application/pdf" // Content-Type advertised for a real sheet
)

// looksLikePDF reports whether a response body is a real rate-sheet PDF.
//
// The new /assets/document host serves a 200 with a non-PDF body for missing
// sheets instead of a 404, so the status code alone is not enough to tell a
// real rate sheet from a "not found" response.
//
// The body's %PDF signature is authoritative: the Content-Type header is not
// always trustworthy (per PR review), so we never accept on the header alone.
// It is used only to corroborate — a body that passes the signature check but
// whose Content-Type explicitly advertises a non-PDF type is still rejected,
// which catches an HTML soft-404 even if its body were to start with %PDF.
func looksLikePDF(contentType string, content []byte) bool {
	if !bytes.HasPrefix(content, []byte(pdfSignature)) {
		return false
	}

	ct := strings.ToLower(contentType)
	if ct != "" && !strings.Contains(ct, pdfContentType) {
		return false
	}

	return true
}

// SBP has hosted the daily sheets under three naming schemes as it migrated to
// the /assets/document host. ratePath picks the right one for a date:
//   - migrationOverrides: transition-window sheets with irregular names
//   - after 2026-07-02: current prefix + DD-month-YYYY (e.g. 03-july-2026)
//   - after 2026-05-31: bare DD-Mon-YY (e.g. 23-Jun-26)
//   - earlier: prefix + DD-Mon-YY (e.g. 27-Aug-25)
//
// The boundaries are the last day of the previous scheme so the checks can use
// date.After; dates are day-truncated by ForDate/ForTime.
var (
	lastLegacyDay   = time.Date(2026, time.July, 2, 0, 0, 0, 0, time.UTC) // current scheme starts the next day
	lastPrefixedDay = time.Date(2026, time.May, 31, 0, 0, 0, 0, time.UTC) // bare scheme starts the next day
)

// migrationOverrides holds the transition-window sheets that SBP uploaded with
// irregular names that fit none of the date-based schemes. Keyed by YYYY-MM-DD.
var migrationOverrides = map[string]string{
	"2026-06-30": "/30-Jun-26_1.pdf",            // bare date with a _1 suffix
	"2026-07-02": ratePrefix + "-02-Jul-26.pdf", // prefix but legacy date style
}

// ratePath builds the URL path for the exchange rate PDF of the given date.
func ratePath(date time.Time) string {
	if override, ok := migrationOverrides[date.Format("2006-01-02")]; ok {
		return override
	}

	// Legacy dates: DD-Mon-YY, e.g. 27-Aug-25.
	legacyDate := fmt.Sprintf("%02d-%s-%02d",
		date.Day(),
		date.Format("Jan"),
		date.Year()%YearModulo, // Last 2 digits of year
	)

	switch {
	case date.After(lastLegacyDay):
		// Current scheme: prefix + DD-month-YYYY, e.g. 03-july-2026.
		return fmt.Sprintf("%s-%02d-%s-%d.pdf",
			ratePrefix,
			date.Day(),
			strings.ToLower(date.Format("January")),
			date.Year(),
		)
	case date.After(lastPrefixedDay):
		// Recent legacy window: bare name, e.g. /23-Jun-26.pdf.
		return "/" + legacyDate + ".pdf"
	default:
		// Older archive: prefix + DD-Mon-YY, e.g. .../rate-27-Aug-25.pdf.
		return ratePrefix + "-" + legacyDate + ".pdf"
	}
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

	if !looksLikePDF(resp.Header.Get("Content-Type"), content) {
		return nil, fmt.Errorf("PDF not found: no rate sheet available for path: %s", path)
	}

	rates, err := parsePDFContent(content, date, fullURL)
	if err != nil {
		// The PDF exists but isn't a parseable rate sheet. SBP posted a few
		// malformed/unrelated PDFs during the June 2026 migration (e.g.
		// 2026-06-01, 03, 04, 05); surface a clear, date-tagged error rather
		// than the raw parser message so callers can distinguish it from a bug.
		return nil, fmt.Errorf("no valid rate sheet for %s (%s): %w", date.Format("2006-01-02"), fullURL, err)
	}

	return rates, nil
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

	if !looksLikePDF(resp.Header.Get("Content-Type"), content) {
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
