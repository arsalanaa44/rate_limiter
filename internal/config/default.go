package config

import (
	"github.com/arsalanaa44/rate_limiter/internal/db/redis"
)

func Default() Config {
	return Config{
		Debug: true,
		Port:  8080,
		Database: redis.Config{
			URL: "0.0.0.0:6379",
		},
	}
}
