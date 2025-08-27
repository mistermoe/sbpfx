---
title: State Bank of Pakistan Exchange Rates
description: Fetch exchange rates from the State Bank of Pakistan
---

The State Bank of Pakistan (SBP) publishes exchange rates for various currencies against the Pakistani Rupee (PKR) on a daily basis. They publish the exchange rates on the internet in PDF format at a deterministic URL for each day. This lib provides an interface to fetch and/or download the exchange rates for any given date.

Example Rate Sheets:
* [2025-08-27](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/27-Aug-25.pdf)
* [2025-08-26](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/26-Aug-25.pdf)
* [2025-08-25](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/25-Aug-25.pdf)
* [2025-08-24](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/24-Aug-25.pdf)
* [2025-08-23](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/23-Aug-25.pdf)
* [2025-08-22](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/22-Aug-25.pdf)


:::note
These Exchange Rates are issued by the State Bank of Pakistan for Authorized Dealers to revalue their books daily on Mark-to-Market basis. M2M rate of USD is compiled as weighted average of closing interbank exchange rate collected through Brokerage Houses. M2M rates of other currencies are compiled on the basis of USD/PKR rate compiled from brokerage houses' data and exchange rate of other currencies against USD quoted on Reuters Eikon Terminal.
:::


:::caution[Important]
These rates will almost always differ from google rates (e.g. rates from google search).
:::



## Installation

```bash
go get github.com/mistermoe/sbpfx
```

## Usage

```go
import "github.com/mistermoe/sbpfx"

func main() {
	client := sbpfx.New()
	rate, err := client.GetExchangeRate(context.Background(), sbpfx.USD)
	if err != nil {
		log.Fatalf("failed to get exchange rate: %v", err)
	}
	fmt.Println(rate)
}
```

## API

### `GetExchangeRate`

Fetches the exchange rate for a given currency. Optionally, you can pass in a date to fetch the exchange rate for a specific date. If no date is provided, the current date is used.

```go
import "github.com/mistermoe/sbpfx"

func main() {
	client := sbpfx.New()
	rate, err := client.GetExchangeRate(context.Background(), sbpfx.USD)
	if err != nil {
		log.Fatalf("failed to get exchange rate: %v", err)
	}
	
  rateJSON, err := json.Marshal(rate)
  if err != nil {
    log.Fatalf("failed to marshal exchange rate: %v", err)
  }
  
  fmt.Println(string(rateJSON))
  
  // with date
  rate, err = client.GetExchangeRate(context.Background(), sbpfx.USD, sbpfx.ForDate("2025-08-27"))
  if err != nil {
    log.Fatalf("failed to get exchange rate: %v", err)
  }
  fmt.Println(rate)
}
```


### `GetExchangeRates`

Fetches the exchange rates for all currencies. Optionally, you can pass in a date to fetch the exchange rates for a specific date. If no date is provided, the current date is used.

```go
import "github.com/mistermoe/sbpfx"

func main() {
	client := sbpfx.New()
	rates, err := client.GetExchangeRates(context.Background())
	if err != nil {
		log.Fatalf("failed to get exchange rates: %v", err)
	}
	
	ratesJSON, err := json.Marshal(rates)
	if err != nil {
		log.Fatalf("failed to marshal exchange rates: %v", err)
	}
	
	fmt.Println(string(ratesJSON))
	
	// with date
	rates, err = client.GetExchangeRates(context.Background(), sbpfx.ForDate("2025-08-27"))
	if err != nil {
		log.Fatalf("failed to get exchange rates: %v", err)
	}
	
	ratesJSON, err = json.Marshal(rates)
	if err != nil {
		log.Fatalf("failed to marshal exchange rates: %v", err)
	}
	
	fmt.Println(string(ratesJSON))
}
```

### `DownloadRateSheet`

Downloads the exchange rate sheet for a given date to the specified file path. Optionally, you can pass in a date to download the exchange rate sheet for a specific date. If no date is provided, the current date is used.

```go
import "github.com/mistermoe/sbpfx"

func main() {
	client := sbpfx.New()
	err := client.DownloadRateSheet(context.Background(), "exchange_rates.pdf")
	if err != nil {
		log.Fatalf("failed to download exchange rate sheet: %v", err)
	}
	
	// with date
	err = client.DownloadRateSheet(context.Background(), "exchange_rates.pdf", sbpfx.ForDate("2025-08-27"))
	if err != nil {
		log.Fatalf("failed to download exchange rate sheet: %v", err)
	}
}
```


### `GetUrl`

Returns the URL for the exchange rate sheet for a given date. Optionally, you can pass in a date to get the URL for a specific date. If no date is provided, the current date is used.

```go
import "github.com/mistermoe/sbpfx"

func main() {
	client := sbpfx.New()
	url := client.GetUrl()
	fmt.Println(url)
	
	// with date
	url = client.GetUrl(sbpfx.ForDate("2025-08-27"))
	fmt.Println(url)
}
```