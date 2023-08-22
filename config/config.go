package config

import (
	"os"
)

type Config struct {
	PolygonKey string
}

func GetConfig() Config {
	c := Config{
		PolygonKey: os.Getenv("POLYGON_KEY"),
	}
	return c
}
