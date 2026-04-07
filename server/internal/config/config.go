package config

import (
	"os"
	"strconv"
)

type Config struct {
	Latitude      float64
	Longitude     float64
	ThresholdMPH  float64
	ExpoPushToken string
	DBPath        string
}

func Load() *Config {
	lat, _ := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	lon, _ := strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)

	threshold, _ := strconv.ParseFloat(os.Getenv("WIND_THRESHOLD_MPH"), 64)
	if threshold == 0 {
		threshold = 20
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "breeze.db"
	}

	return &Config{
		Latitude:      lat,
		Longitude:     lon,
		ThresholdMPH:  threshold,
		ExpoPushToken: os.Getenv("EXPO_PUSH_TOKEN"),
		DBPath:        dbPath,
	}
}
