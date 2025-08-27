# Improved Developer Experience with ForDate and ForTime

## Before (time.Time only):
```go
// Had to construct time.Time manually
specificDate := time.Date(2025, 8, 27, 0, 0, 0, 0, time.UTC)
rate, err := client.GetExchangeRate(ctx, USD, ForDate(specificDate))
```

## After (human-friendly string format):
```go
// Much cleaner and more readable
rate, err := client.GetExchangeRate(ctx, USD, ForDate("2025-08-27"))

// Still supports time.Time for programmatic use
specificTime := time.Date(2025, 8, 27, 0, 0, 0, 0, time.UTC)
rate, err := client.GetExchangeRate(ctx, USD, ForTime(specificTime))
```

## Usage Examples:

### String-based (recommended for most cases):
```go
client := sbp.New()
ctx := context.Background()

// Today's rates (default)
rate, err := client.GetExchangeRate(ctx, sbp.USD)

// Specific date with string
rate, err := client.GetExchangeRate(ctx, sbp.USD, sbp.ForDate("2025-08-27"))

// Multiple options still work
url := client.GetUrl(sbp.ForDate("2025-12-25"))
err := client.DownloadRateSheet(ctx, "rates.pdf", sbp.ForDate("2025-08-27"))
```

### Time-based (for programmatic use):
```go
// When working with existing time.Time values
for _, date := range calculateBusinessDays() {
    rate, err := client.GetExchangeRate(ctx, sbp.USD, sbp.ForTime(date))
    // process rates...
}
```

## Benefits:
- ✅ Human-readable date format: "2025-08-27"
- ✅ No manual time.Time construction needed
- ✅ Backward compatibility with ForTime()
- ✅ Works with all client methods
- ✅ Automatic UTC conversion and day truncation
- ✅ Graceful error handling for invalid formats
