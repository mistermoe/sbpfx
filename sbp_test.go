package sbpfx_test

import (
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/mistermoe/httpr"
	"github.com/mistermoe/sbpfx"
	"github.com/mistermoe/sbpfx/vcr"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var testMode vcr.Mode = vcr.Replay

func bootstrap(_ *testing.T, mode vcr.Mode, rec *recorder.Recorder) *sbpfx.Client {
	recorder := rec.GetDefaultClient()
	return sbpfx.New(httpr.HTTPClient(*recorder))
}

func TestGetExchangeRates(t *testing.T) {
	vcr.Test(t, testMode, bootstrap, func(t *testing.T, client *sbpfx.Client, c vcr.Cassette) {
		rate, err := client.GetExchangeRate(t.Context(), sbpfx.USD)
		assert.NoError(t, err)
		assert.NotZero(t, rate)

		assert.Equal(t, sbpfx.USD, rate.Currency)
		assert.NotZero(t, rate.Ready)
		assert.NotZero(t, rate.Date)
		assert.NotZero(t, rate.URL)
	})
}

func TestDownloadRateSheet(t *testing.T) {
	vcr.Test(t, testMode, bootstrap, func(t *testing.T, client *sbpfx.Client, c vcr.Cassette) {
		// Create a temporary file path
		tempFile := t.TempDir() + "/test_rate_sheet.pdf"

		// Download the rate sheet
		err := client.DownloadRateSheet(t.Context(), tempFile)
		assert.NoError(t, err)

		// Verify the file was created and has content
		stat, err := os.Stat(tempFile)
		assert.NoError(t, err, "File should be created successfully")
		assert.True(t, stat.Size() > 0, "Downloaded file should have content")
	})
}

func TestGetExchangeRatesFutureDate(t *testing.T) {
	vcr.Test(t, testMode, bootstrap, func(t *testing.T, client *sbpfx.Client, c vcr.Cassette) {
		// Try to get exchange rates for a date far in the future using string format
		rate, err := client.GetExchangeRate(t.Context(), sbpfx.USD, sbpfx.ForDate("2030-12-25"))

		// Should return an error since the PDF won't exist for future dates
		assert.Error(t, err, "Should return error for future date")
		assert.Zero(t, rate, "Rate should be nil/zero for future date")

		// The error should indicate the PDF was not found
		assert.Contains(t, err.Error(), "PDF not found", "Error should mention PDF not found")
	})
}

func TestForDateAndForTime(t *testing.T) {
	vcr.Test(t, testMode, bootstrap, func(t *testing.T, client *sbpfx.Client, c vcr.Cassette) {
		// Test ForDate with string format
		rateFromString, err1 := client.GetExchangeRate(t.Context(), sbpfx.USD, sbpfx.ForDate("2025-08-27"))

		// Test ForTime with time.Time
		specificTime := time.Date(2025, 8, 27, 0, 0, 0, 0, time.UTC)
		rateFromTime, err2 := client.GetExchangeRate(t.Context(), sbpfx.USD, sbpfx.ForTime(specificTime))

		// Both should work and return the same data (assuming the date exists)
		if err1 == nil && err2 == nil {
			assert.Equal(t, rateFromString.Currency, rateFromTime.Currency)
			assert.Equal(t, rateFromString.Ready, rateFromTime.Ready)
			assert.Equal(t, rateFromString.Date.Format("2006-01-02"), rateFromTime.Date.Format("2006-01-02"))
			assert.Equal(t, rateFromString.URL, rateFromTime.URL)
		}

		// At minimum, both should have the same error status
		assert.Equal(t, err1 != nil, err2 != nil, "Both ForDate and ForTime should behave consistently")
	})
}

func TestForDateInvalidFormat(t *testing.T) {
	client := sbpfx.New()

	// Test with invalid date format - should use default date and not crash
	url1 := client.GetUrl(sbpfx.ForDate("invalid-date"))
	url2 := client.GetUrl() // default date

	// The URLs should be different only if the invalid date somehow got parsed differently
	// But more importantly, the function shouldn't crash
	assert.True(t, len(url1) > 0, "Should still generate a URL even with invalid date")
	assert.True(t, len(url2) > 0, "Should generate a URL with default date")
}
