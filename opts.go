package sbpfx

import (
	"fmt"
	"time"
)

const (
	HoursInDay = 24 // Hours in a day for time truncation
)

type Option func(*option) error

type option struct {
	date time.Time
}

// ForDate sets a specific date for the exchange rate request using a string in YYYY-MM-DD format.
func ForDate(dateStr string) Option {
	return func(c *option) error {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("invalid date format '%s', expected format: YYYY-MM-DD", dateStr)
		}

		c.date = parsedDate.UTC().Truncate(HoursInDay * time.Hour)
		return nil
	}
}

// ForTime sets a specific date for the exchange rate request using a time.Time.
func ForTime(date time.Time) Option {
	return func(c *option) error {
		c.date = date.UTC().Truncate(HoursInDay * time.Hour)
		return nil
	}
}

func defaultConfig() *option {
	return &option{
		date: time.Now().UTC().Truncate(HoursInDay * time.Hour),
	}
}
