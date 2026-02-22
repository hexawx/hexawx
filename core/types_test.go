package core

import (
	"testing"
	"time"
)

func TestWeatherRecord(t *testing.T) {
	now := time.Now()
	record := WeatherRecord{
		Temperature: 25.5,
		Humidity:    60.0,
		Timestamp:   now,
	}

	if record.Temperature != 25.5 {
		t.Errorf("Attendu 25.5, obtenu %f", record.Temperature)
	}
}
