package config

import (
	"os"
	"time"
)

func GetSecrets() Secrets {
	s := Secrets{
		PolygonKey: os.Getenv("POLYGON_KEY"),
	}
	return s
}

func GetConfig() Config {
	c := Config{
		Tickers:        []string{"AAPL", "MSFT", "GOOG", "AMZN"},
		NumberOfJobs:   5,
		JobsStartDate:  time.Date(2021, time.August, 1, 0, 0, 0, 0, time.UTC),
		JobsEndDate:    time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
		JobRangeLength: 1,
		JobRangeType:   "minute",
		JobAdjusted:    true,
		JobLimit:       5000,
	}
	return c
}

type Secrets struct {
	PolygonKey string
}

type Config struct {
	Tickers        []string
	NumberOfJobs   int
	JobRangeLength int
	JobRangeType   string
	JobsStartDate  time.Time
	JobsEndDate    time.Time
	JobAdjusted    bool
	JobLimit       int
}
