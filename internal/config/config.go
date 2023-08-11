package config

import (
	"github.com/arsalanaa44/rate_limiter/internal/db/redis"
)

type Config struct {
	Debug    bool         `koanf:"debug"`
	Port     int          `koanf:"port"`
	Database redis.Config `koanf:"database"`
}
