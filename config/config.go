package config

import (
	"os"
)

func GetSecrets() Secrets {
	s := Secrets{
		PolygonKey: os.Getenv("POLYGON_KEY"),
	}
	return s
}

func GetConfig() Config {
	c := Config{
		Tickers:       []string{"AAPL", "MSFT", "GOOG", "AMZN"},
		NumberOfJobs:  5,
		JobsStartDate: "2021-08-01",
		JobsEndDate:   "2023-08-01",
	}
	return c
}

type Secrets struct {
	PolygonKey string
}

type Config struct {
	Tickers       []string
	NumberOfJobs  int
	JobsStartDate string
	JobsEndDate   string
}
