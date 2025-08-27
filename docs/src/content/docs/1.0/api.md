---
title: API Reference
description: Complete API reference for the sbpfx library
slug: 1.0/api
---

# API Reference

The sbpfx library provides a simple, clean API for fetching Pakistani Rupee exchange rates from the State Bank of Pakistan.

## Client

### `New(options ...httpr.ClientOption) *Client`

Creates a new client instance with optional HTTP client configuration.

```go
// Basic client
client := sbpfx.New()

// Client with custom timeout
client := sbpfx.New(httpr.Timeout(60*time.Second))
```

## Exchange Rate Methods

### `GetExchangeRate(ctx context.Context, currency Currency, opts ...Option) (*ExchangeRate, error)`

Fetches the exchange rate for a specific currency.

```go
// Get today's USD rate
rate, err := client.GetExchangeRate(ctx, sbpfx.USD)

// Get USD rate for specific date
rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForDate("2025-08-27"))
```

**Parameters:**

* `ctx`: Context for request cancellation and timeouts
* `currency`: Currency code (e.g., `sbpfx.USD`, `sbpfx.EUR`)
* `opts`: Optional date specification

**Returns:**

* `*ExchangeRate`: Exchange rate data
* `error`: Error if the request fails

### `GetExchangeRates(ctx context.Context, opts ...Option) (map[Currency]*ExchangeRate, error)`

Fetches exchange rates for all available currencies.

```go
// Get today's rates for all currencies
rates, err := client.GetExchangeRates(ctx)

// Get rates for specific date
rates, err := client.GetExchangeRates(ctx, sbpfx.ForDate("2025-08-27"))

// Access specific currency
usdRate := rates[sbpfx.USD]
```

**Parameters:**

* `ctx`: Context for request cancellation and timeouts
* `opts`: Optional date specification

**Returns:**

* `map[Currency]*ExchangeRate`: Map of currency codes to exchange rate data
* `error`: Error if the request fails

## Utility Methods

### `GetUrl(opts ...Option) string`

Returns the URL for the exchange rate PDF for a given date.

```go
// Get today's PDF URL
url := client.GetUrl()

// Get PDF URL for specific date
url := client.GetUrl(sbpfx.ForDate("2025-08-27"))
```

### `DownloadRateSheet(ctx context.Context, path string, opts ...Option) error`

Downloads the original PDF rate sheet to a file.

```go
// Download today's rate sheet
err := client.DownloadRateSheet(ctx, "rates.pdf")

// Download rate sheet for specific date
err := client.DownloadRateSheet(ctx, "rates.pdf", sbpfx.ForDate("2025-08-27"))
```

**Parameters:**

* `ctx`: Context for request cancellation and timeouts
* `path`: Local file path where PDF will be saved
* `opts`: Optional date specification

## Options

### `ForDate(dateStr string) Option`

Specifies a date using a human-readable string format.

```go
// Use YYYY-MM-DD format
rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForDate("2025-08-27"))
```

### `ForTime(date time.Time) Option`

Specifies a date using a `time.Time` value.

```go
// Use time.Time for programmatic date handling
specificTime := time.Date(2025, 8, 27, 0, 0, 0, 0, time.UTC)
rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForTime(specificTime))
```

## Data Types

### `Currency`

String-based currency type with predefined constants.

**Available Currencies:**

* `USD`, `EUR`, `GBP`, `JPY`, `CHF`
* `AUD`, `CAD`, `SEK`, `NOK`, `DKK`
* `SAR`, `AED`, `KWD`, `BHD`, `QAR`, `OMR`
* `CNY`, `HKD`, `SGD`, `THB`, `MYR`, `INR`, `KRW`
* `NZD`, `ZAR`, `BDT`, `BRL`, `ARS`, `LKR`, `TRY`, `IDR`, `MXN`, `RUB`, `GNH`

### `ExchangeRate`

Contains exchange rate data for a specific currency and date.

```go
type ExchangeRate struct {
    Currency Currency  `json:"currency"`
    Date     time.Time `json:"date"`
    URL      string    `json:"url"`        // Source PDF URL
    
    // Spot and Forward Rates (all against PKR)
    Ready      string `json:"ready,omitempty"`       // Spot rate
    OneWeek    string `json:"one_week,omitempty"`    // 1-week forward
    TwoWeek    string `json:"two_week,omitempty"`    // 2-week forward
    OneMonth   string `json:"one_month,omitempty"`   // 1-month forward
    TwoMonth   string `json:"two_month,omitempty"`   // 2-month forward
    ThreeMonth string `json:"three_month,omitempty"` // 3-month forward
    FourMonth  string `json:"four_month,omitempty"`  // 4-month forward
    FiveMonth  string `json:"five_month,omitempty"`  // 5-month forward
    SixMonth   string `json:"six_month,omitempty"`   // 6-month forward
    NineMonth  string `json:"nine_month,omitempty"`  // 9-month forward
    OneYear    string `json:"one_year,omitempty"`    // 1-year forward
}
```

**Methods:**

* `GetSpotRate() string`: Returns the spot rate (Ready rate) as a string

## Error Handling

The library returns standard Go errors. Common error scenarios:

* **Network errors**: Connection failures, timeouts
* **HTTP errors**: 404 for non-existent dates, server errors
* **Parsing errors**: PDF format changes or corruption
* **File errors**: Permission issues when downloading PDFs

```go
rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForDate("2030-12-25"))
if err != nil {
    if strings.Contains(err.Error(), "PDF not found") {
        // Handle missing date (weekends, holidays, future dates)
        log.Printf("No rates available for the requested date")
    } else {
        // Handle other errors
        log.Printf("Error fetching rates: %v", err)
    }
}
```
