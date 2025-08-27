---
title: State Bank of Pakistan Exchange Rates
description: Fetch exchange rates from the State Bank of Pakistan
slug: "1.0"
---

# State Bank of Pakistan Exchange Rates

[![Test](https://github.com/mistermoe/sbpfx/actions/workflows/test.yml/badge.svg)](https://github.com/mistermoe/sbpfx/actions/workflows/test.yml)
[![Lint](https://github.com/mistermoe/sbpfx/actions/workflows/lint.yml/badge.svg)](https://github.com/mistermoe/sbpfx/actions/workflows/lint.yml)

## Overview

The State Bank of Pakistan (SBP) publishes exchange rates for various currencies against the Pakistani Rupee (PKR) on a daily basis. They publish the exchange rates on the internet in PDF format at a deterministic URL for each day. This library provides an interface to fetch and/or download the exchange rates for any given date.

### Example Rate Sheets

* [2025-08-27](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/27-Aug-25.pdf)
* [2025-08-26](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/26-Aug-25.pdf)
* [2025-08-25](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/25-Aug-25.pdf)
* [2025-08-24](https://www.sbp.org.pk/ecodata/rates/m2m/2025/Aug/24-Aug-25.pdf)

:::note
These Exchange Rates are issued by the State Bank of Pakistan for Authorized Dealers to revalue their books daily on Mark-to-Market basis. M2M rate of USD is compiled as weighted average of closing interbank exchange rate collected through Brokerage Houses. M2M rates of other currencies are compiled on the basis of USD/PKR rate compiled from brokerage houses' data and exchange rate of other currencies against USD quoted on Reuters Eikon Terminal.
:::

:::caution\[Important]
These rates will almost always differ from google rates (e.g. rates from google search).
:::

## Installation

```bash
go get github.com/mistermoe/sbpfx
```

## Quick Start

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
        log.Fatalf("failed to get exchange rate: %v", err)
    }
    
    fmt.Printf("USD to PKR: %s\n", rate.GetSpotRate())
}
```

## Features

* **📈 Current Exchange Rates**: Get today's rates with no configuration
* **📅 Historical Data**: Fetch rates for any specific date
* **💾 PDF Download**: Download original PDF rate sheets
* **🔗 URL Generation**: Get direct links to rate sheet PDFs
* **🌍 Multiple Currencies**: Support for USD, EUR, GBP, JPY, and 25+ other currencies
* **⚡ Human-Friendly API**: Use simple date strings like "2025-08-27"
* **🧪 Well Tested**: Comprehensive test suite with VCR for reliable testing
* **📦 Zero Dependencies**: Minimal external dependencies for easy integration
