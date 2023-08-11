package redis

import (
	"time"
)

type Config struct {
	URL        string        `koanf:"url"`
	Expiration time.Duration `koanf:"expiration"`
}
