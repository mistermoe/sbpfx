---
title: Examples
description: Practical examples of using the sbpfx library
slug: 1.0/examples
---

# Examples

This page provides practical examples of using the sbpfx library for various use cases.

## Basic Usage

### Get Today's Exchange Rate

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    rate, err := client.GetExchangeRate(context.Background(), sbpfx.USD)
    if err != nil {
        log.Fatalf("failed to get USD rate: %v", err)
    }
    
    fmt.Printf("USD to PKR rate: %s\n", rate.GetSpotRate())
    fmt.Printf("Date: %s\n", rate.Date.Format("2006-01-02"))
    fmt.Printf("Source: %s\n", rate.URL)
}
```

### Get All Exchange Rates

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    rates, err := client.GetExchangeRates(context.Background())
    if err != nil {
        log.Fatalf("failed to get exchange rates: %v", err)
    }
    
    // Print as JSON
    ratesJSON, err := json.MarshalIndent(rates, "", "  ")
    if err != nil {
        log.Fatalf("failed to marshal rates: %v", err)
    }
    
    fmt.Println(string(ratesJSON))
}
```

## Working with Dates

### Using String Dates (Recommended)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    // Get rate for specific date using human-readable format
    rate, err := client.GetExchangeRate(
        context.Background(), 
        sbpfx.USD, 
        sbpfx.ForDate("2025-08-27"),
    )
    if err != nil {
        log.Fatalf("failed to get USD rate: %v", err)
    }
    
    fmt.Printf("USD rate on 2025-08-27: %s\n", rate.GetSpotRate())
}
```

### Using time.Time

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    // Get rates for a specific time.Time
    specificDate := time.Date(2025, 8, 27, 0, 0, 0, 0, time.UTC)
    rate, err := client.GetExchangeRate(
        context.Background(), 
        sbpfx.USD, 
        sbpfx.ForTime(specificDate),
    )
    if err != nil {
        log.Fatalf("failed to get USD rate: %v", err)
    }
    
    fmt.Printf("USD rate on %s: %s\n", 
        specificDate.Format("2006-01-02"), 
        rate.GetSpotRate(),
    )
}
```

## PDF Operations

### Download Rate Sheet

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    // Download today's rate sheet
    err := client.DownloadRateSheet(context.Background(), "today_rates.pdf")
    if err != nil {
        log.Fatalf("failed to download rate sheet: %v", err)
    }
    
    fmt.Println("Rate sheet downloaded successfully!")
}
```

### Get PDF URL

```go
package main

import (
    "fmt"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    // Get URL for today's rate sheet
    url := client.GetUrl()
    fmt.Printf("Today's rate sheet URL: %s\n", url)
    
    // Get URL for specific date
    url = client.GetUrl(sbpfx.ForDate("2025-08-27"))
    fmt.Printf("Rate sheet URL for 2025-08-27: %s\n", url)
}
```

## Advanced Examples

### Historical Rate Analysis

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strconv"
    "time"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    ctx := context.Background()
    
    // Analyze USD rates for the past week
    fmt.Println("USD to PKR rates for the past week:")
    
    for i := 0; i < 7; i++ {
        date := time.Now().AddDate(0, 0, -i)
        dateStr := date.Format("2006-01-02")
        
        rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForTime(date))
        if err != nil {
            fmt.Printf("%s: Error - %v\n", dateStr, err)
            continue
        }
        
        // Convert rate to float for analysis
        rateFloat, err := strconv.ParseFloat(rate.GetSpotRate(), 64)
        if err != nil {
            fmt.Printf("%s: Invalid rate format\n", dateStr)
            continue
        }
        
        fmt.Printf("%s: %.4f PKR\n", dateStr, rateFloat)
    }
}
```

### Multiple Currencies Comparison

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    ctx := context.Background()
    
    // Compare major currencies
    currencies := []sbpfx.Currency{
        sbpfx.USD, sbpfx.EUR, sbpfx.GBP, 
        sbpfx.JPY, sbpfx.CHF, sbpfx.AUD,
    }
    
    fmt.Println("Major currency rates against PKR:")
    
    for _, currency := range currencies {
        rate, err := client.GetExchangeRate(ctx, currency)
        if err != nil {
            fmt.Printf("%-3s: Error - %v\n", currency, err)
            continue
        }
        
        fmt.Printf("%-3s: %s PKR\n", currency, rate.GetSpotRate())
    }
}
```

### Error Handling Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    ctx := context.Background()
    
    // Try to get rate for a future date (will fail)
    rate, err := client.GetExchangeRate(ctx, sbpfx.USD, sbpfx.ForDate("2030-12-25"))
    if err != nil {
        if strings.Contains(err.Error(), "PDF not found") {
            fmt.Println("No rates available for future dates")
        } else if strings.Contains(err.Error(), "context deadline exceeded") {
            fmt.Println("Request timed out")
        } else {
            fmt.Printf("Unexpected error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("USD rate: %s\n", rate.GetSpotRate())
}
```

### With Custom HTTP Configuration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mistermoe/httpr"
    "github.com/mistermoe/sbpfx"
)

func main() {
    // Create client with custom timeout and user agent
    client := sbpfx.New(
        httpr.Timeout(60*time.Second),
        httpr.UserAgent("MyApp/1.0"),
    )
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    rate, err := client.GetExchangeRate(ctx, sbpfx.USD)
    if err != nil {
        log.Fatalf("failed to get USD rate: %v", err)
    }
    
    fmt.Printf("USD rate: %s\n", rate.GetSpotRate())
}
```

## JSON Output Examples

### Single Rate as JSON

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    rate, err := client.GetExchangeRate(context.Background(), sbpfx.USD)
    if err != nil {
        log.Fatalf("failed to get exchange rate: %v", err)
    }
    
    rateJSON, err := json.MarshalIndent(rate, "", "  ")
    if err != nil {
        log.Fatalf("failed to marshal exchange rate: %v", err)
    }
    
    fmt.Println(string(rateJSON))
}
```

### All Rates as JSON

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/mistermoe/sbpfx"
)

func main() {
    client := sbpfx.New()
    
    rates, err := client.GetExchangeRates(context.Background())
    if err != nil {
        log.Fatalf("failed to get exchange rates: %v", err)
    }
    
    ratesJSON, err := json.MarshalIndent(rates, "", "  ")
    if err != nil {
        log.Fatalf("failed to marshal exchange rates: %v", err)
    }
    
    fmt.Println(string(ratesJSON))
}
```
