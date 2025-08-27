---
title: Troubleshooting
description: Common issues and solutions when using the sbpfx library
slug: 1.0/troubleshooting
---

# Troubleshooting

This page covers common issues and their solutions when using the sbpfx library.

## Common Issues

### Rate Not Found for Specific Date

**Problem:** Getting "PDF not found" error for certain dates.

**Possible Causes:**

* **Weekends**: SBP doesn't publish rates on Saturdays and Sundays
* **Holidays**: No rates published on Pakistani public holidays
* **Future dates**: Rates don't exist for dates that haven't occurred yet
* **Very old dates**: Historical data might not be available

**Solution:**

```go
rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForDate("2025-08-27"))
if err != nil {
    if strings.Contains(err.Error(), "PDF not found") {
        // Try previous business day or handle missing data
        fmt.Println("No rates available for this date")
    }
}
```

### Timeout Errors

**Problem:** Getting "context deadline exceeded" or timeout errors.

**Possible Causes:**

* Slow network connection
* SBP server issues
* Large PDF files taking time to download

**Solutions:**

1. **Increase timeout:**

```go
// Increase client timeout
client := sbpfx.New(httpr.Timeout(60*time.Second))

// Or use context timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

2. **Retry logic:**

```go
func getExchangeRateWithRetry(client *sbpfx.Client, currency sbpfx.Currency, maxRetries int) (*sbpfx.ExchangeRate, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        rate, err := client.GetExchangeRate(context.Background(), currency)
        if err == nil {
            return rate, nil
        }
        lastErr = err
        time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
    }
    return nil, lastErr
}
```

### Invalid Date Format

**Problem:** Date string not being parsed correctly.

**Solution:** Ensure you're using the correct format:

```go
// ✅ Correct format
rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForDate("2025-08-27"))

// ❌ Incorrect formats
// ForDate("08/27/2025")  // Wrong format
// ForDate("27-08-2025")  // Wrong format
// ForDate("2025-8-27")   // Missing zero padding
```

### PDF Parsing Errors

**Problem:** Error parsing PDF content or unexpected data structure.

**Possible Causes:**

* SBP changed PDF format
* Corrupted PDF download
* Network issues during download

**Solution:**

```go
// Download PDF manually for inspection
err := client.DownloadRateSheet(ctx, "debug.pdf", sbpfx.ForDate("2025-08-27"))
if err != nil {
    log.Printf("Failed to download PDF: %v", err)
    // PDF might be corrupted or format changed
}
```

## Debugging Tips

### Enable HTTP Debugging

If you're using httpr's debug features:

```go
import "github.com/mistermoe/httpr"

client := sbpfx.New(
    httpr.Inspect(), // Logs HTTP requests and responses
)
```

### Check PDF URLs

Verify the URL generation is correct:

```go
url := client.GetUrl(sbpfx.ForDate("2025-08-27"))
fmt.Printf("Generated URL: %s\n", url)

// Manually check if URL is accessible in browser
```

### Validate Exchange Rate Data

```go
rate, err := client.GetExchangeRate(ctx, sbpfx.USD)
if err == nil {
    fmt.Printf("Currency: %s\n", rate.Currency)
    fmt.Printf("Date: %s\n", rate.Date.Format("2006-01-02"))
    fmt.Printf("Spot Rate: %s\n", rate.GetSpotRate())
    fmt.Printf("Source URL: %s\n", rate.URL)
    
    // Validate rate is reasonable
    if rateFloat, err := strconv.ParseFloat(rate.GetSpotRate(), 64); err == nil {
        if rateFloat < 100 || rateFloat > 500 {
            fmt.Printf("Warning: Rate seems unusual: %.2f\n", rateFloat)
        }
    }
}
```

## FAQ

### Q: Why do rates differ from Google or other sources?

**A:** SBP rates are official interbank rates used by authorized dealers for mark-to-market revaluation. They differ from consumer exchange rates shown by Google, which typically include retail margins.

### Q: What time are new rates published?

**A:** SBP typically publishes rates during Pakistani business hours. Rates for a given date are usually available by end of business day in Pakistan time (PKT).

### Q: Are rates available for weekends?

**A:** No, SBP doesn't publish rates on weekends (Saturday/Sunday) or public holidays in Pakistan.

### Q: How far back do historical rates go?

**A:** The availability of historical data depends on SBP's archive. Very old dates might not be available.

### Q: Can I get intraday rates?

**A:** No, SBP publishes daily rates only. For real-time rates, you would need a different data source.

### Q: What currencies are supported?

**A:** The library supports 30+ currencies including major currencies (USD, EUR, GBP, JPY) and regional currencies. See the API reference for the complete list.

### Q: How do I handle rate unavailability gracefully?

**A:** Implement fallback logic:

```go
func getAvailableRate(client *sbpfx.Client, currency sbpfx.Currency, preferredDate string) (*sbpfx.ExchangeRate, error) {
    // Try preferred date first
    rate, err := client.GetExchangeRate(ctx, currency, sbpfx.ForDate(preferredDate))
    if err == nil {
        return rate, nil
    }
    
    // Fall back to today's rate
    rate, err = client.GetExchangeRate(ctx, currency)
    if err == nil {
        return rate, nil
    }
    
    return nil, fmt.Errorf("no rates available")
}
```

## Getting Help

If you're still experiencing issues:

1. **Check the GitHub issues**: [github.com/mistermoe/sbpfx/issues](https://github.com/mistermoe/sbpfx/issues)
2. **Create a new issue**: Include error messages, code samples, and environment details
3. **Verify SBP website**: Check if [sbp.org.pk](https://www.sbp.org.pk) is accessible and the PDF format hasn't changed

## Contributing

Found a bug or have a suggestion? Contributions are welcome!

* **Bug reports**: Include steps to reproduce and error messages
* **Feature requests**: Explain the use case and expected behavior
* **Code contributions**: Follow the existing code style and include tests
