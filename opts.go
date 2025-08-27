package sbpfx

import (
	"fmt"
	"time"
)

// Option represents a functional option for configuring exchange rate requests
type Option func(*option)

// option holds the configuration for exchange rate requests
type option struct {
	date time.Time
}

// ForDate sets a specific date for the exchange rate request using a string in YYYY-MM-DD format
func ForDate(dateStr string) Option {
	return func(c *option) {
		// Parse the date string in YYYY-MM-DD format
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			c.date = parsedDate.UTC().Truncate(24 * time.Hour) // Start of day in UTC
		} else {
			// If parsing fails, keep the default date and log the error
			// In a real-world scenario, you might want to handle this differently
			fmt.Printf("Warning: Invalid date format '%s', using default date. Expected format: YYYY-MM-DD\n", dateStr)
		}
	}
}

// ForTime sets a specific date for the exchange rate request using a time.Time
func ForTime(date time.Time) Option {
	return func(c *option) {
		c.date = date.UTC().Truncate(24 * time.Hour) // Start of day in UTC
	}
}

// defaultConfig returns the default configuration (today's date)
func defaultConfig() *option {
	return &option{
		date: time.Now().UTC().Truncate(24 * time.Hour), // Start of day in UTC
	}
}
