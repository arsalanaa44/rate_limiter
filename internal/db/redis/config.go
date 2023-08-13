package redis

type Config struct {
	URL      string `koanf:"url"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}
