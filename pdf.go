package sbpfx

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

// parsePDFContent extracts text from PDF content and parses exchange rates
func parsePDFContent(content []byte, date time.Time, url string) (map[Currency]*ExchangeRate, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var fullText strings.Builder
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		fullText.WriteString(text)
	}

	return parseExchangeRateText(fullText.String(), date, url)
}

// parseExchangeRateText parses extracted text to find exchange rates
func parseExchangeRateText(text string, date time.Time, url string) (map[Currency]*ExchangeRate, error) {
	rates := make(map[Currency]*ExchangeRate)

	lines := strings.Split(text, "\n")

	// Find the CURRENCY header line
	currencyLineIndex := -1
	readyLineIndex := -1

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "CURRENCY" {
			currencyLineIndex = i
		}
		if line == "READY" {
			readyLineIndex = i
		}
		if currencyLineIndex != -1 && readyLineIndex != -1 {
			break
		}
	}

	if currencyLineIndex == -1 || readyLineIndex == -1 {
		return nil, fmt.Errorf("could not find CURRENCY or READY headers")
	}

	// Parse currencies and rates
	currencies := []string{}
	readyRates := []string{}

	// Collect currencies starting after the CURRENCY header
	for i := currencyLineIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || line == "READY" {
			continue
		}
		// Stop when we hit notes or other sections
		if strings.Contains(strings.ToUpper(line), "EXCHANGE RATES FOR MARK") {
			break
		}
		// Check if it's a 3-letter currency code
		if len(line) == 3 && Currency(line).IsValid() {
			currencies = append(currencies, line)
		}
	}

	// Collect READY rates starting after the READY header
	for i := readyLineIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		// Stop when we hit notes or other sections
		if strings.Contains(strings.ToUpper(line), "EXCHANGE RATES FOR MARK") {
			break
		}
		// Check if it's a number (rate)
		if rate, err := strconv.ParseFloat(line, 64); err == nil && rate > 0 {
			readyRates = append(readyRates, line)
		}
	}

	// Match currencies with their ready rates
	minLen := min(len(currencies), len(readyRates))
	for i := 0; i < minLen; i++ {
		// Validate that it's a valid rate string before storing
		if _, err := strconv.ParseFloat(readyRates[i], 64); err == nil {
			currency := Currency(currencies[i])
			rates[currency] = &ExchangeRate{
				Currency: currency,
				Ready:    readyRates[i], // Store as string
				Date:     date,
				URL:      url,
			}
		}
	}

	if len(rates) == 0 {
		return nil, fmt.Errorf("no exchange rates found in PDF")
	}

	return rates, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
