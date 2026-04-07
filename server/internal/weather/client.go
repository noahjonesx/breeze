package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const forecastURL = "https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&hourly=windspeed_10m&wind_speed_unit=mph&forecast_days=2&timezone=auto"

type response struct {
	Hourly struct {
		Time      []string  `json:"time"`
		WindSpeed []float64 `json:"windspeed_10m"`
	} `json:"hourly"`
}

type Entry struct {
	Time      time.Time
	WindSpeed float64
}

func FetchForecast(lat, lon float64) ([]Entry, error) {
	url := fmt.Sprintf(forecastURL, lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("open-meteo returned %d", resp.StatusCode)
	}

	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode forecast: %w", err)
	}

	entries := make([]Entry, 0, len(data.Hourly.Time))
	for i, t := range data.Hourly.Time {
		parsed, err := time.ParseInLocation("2006-01-02T15:04", t, time.Local)
		if err != nil {
			continue
		}
		entries = append(entries, Entry{
			Time:      parsed,
			WindSpeed: data.Hourly.WindSpeed[i],
		})
	}
	return entries, nil
}
