package sbpfx

import "time"

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	CHF Currency = "CHF"
	AUD Currency = "AUD"
	CAD Currency = "CAD"
	SEK Currency = "SEK"
	NOK Currency = "NOK"
	DKK Currency = "DKK"
	SAR Currency = "SAR"
	AED Currency = "AED"
	KWD Currency = "KWD"
	BHD Currency = "BHD"
	QAR Currency = "QAR"
	OMR Currency = "OMR"
	CNY Currency = "CNY"
	HKD Currency = "HKD"
	SGD Currency = "SGD"
	THB Currency = "THB"
	MYR Currency = "MYR"
	INR Currency = "INR"
	KRW Currency = "KRW"
	NZD Currency = "NZD"
	ZAR Currency = "ZAR"
	BDT Currency = "BDT"
	BRL Currency = "BRL"
	ARS Currency = "ARS"
	LKR Currency = "LKR"
	TRY Currency = "TRY"
	IDR Currency = "IDR"
	MXN Currency = "MXN"
	RUB Currency = "RUB"
	GNH Currency = "GNH"
)

func (c Currency) String() string {
	return string(c)
}

func (c Currency) IsValid() bool {
	validCurrencies := map[Currency]bool{
		USD: true, EUR: true, JPY: true, GBP: true, CHF: true,
		AUD: true, CAD: true, SEK: true, NOK: true, DKK: true,
		SAR: true, AED: true, KWD: true, BHD: true, QAR: true,
		OMR: true, CNY: true, HKD: true, SGD: true, THB: true,
		MYR: true, INR: true, KRW: true, NZD: true, ZAR: true,
		BDT: true, BRL: true, ARS: true, LKR: true, TRY: true,
		IDR: true, MXN: true, RUB: true, GNH: true,
	}
	return validCurrencies[c]
}

// ExchangeRate represents exchange rates for different delivery periods
// These are forward rates used for currency hedging and speculation.
type ExchangeRate struct {
	Currency   Currency  `json:"currency"`
	Date       time.Time `json:"date"`
	URL        string    `json:"url"`                   // Source PDF URL
	Ready      string    `json:"ready,omitempty"`       // Spot rate (immediate delivery)
	OneWeek    string    `json:"one_week,omitempty"`    // 1-week forward rate
	TwoWeek    string    `json:"two_week,omitempty"`    // 2-week forward rate
	OneMonth   string    `json:"one_month,omitempty"`   // 1-month forward rate
	TwoMonth   string    `json:"two_month,omitempty"`   // 2-month forward rate
	ThreeMonth string    `json:"three_month,omitempty"` // 3-month forward rate
	FourMonth  string    `json:"four_month,omitempty"`  // 4-month forward rate
	FiveMonth  string    `json:"five_month,omitempty"`  // 5-month forward rate
	SixMonth   string    `json:"six_month,omitempty"`   // 6-month forward rate
	NineMonth  string    `json:"nine_month,omitempty"`  // 9-month forward rate
	OneYear    string    `json:"one_year,omitempty"`    // 1-year forward rate
}

// GetSpotRate returns the spot rate (Ready rate) as a string.
func (e *ExchangeRate) GetSpotRate() string {
	return e.Ready
}
